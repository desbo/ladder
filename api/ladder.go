package app

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"firebase.google.com/go/auth"

	"golang.org/x/net/context"

	"github.com/rs/xid"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

// LadderPlayer is a player in a ladder
type LadderPlayer struct {
	Key      *datastore.Key `json:"key" datastore:"__key__"` // firebase ID
	Name     string         `json:"name"`
	Position int            `json:"position"`
	Wins     int            `json:"wins"`
	Losses   int            `json:"losses"`
	Rating   int            `json:"rating"`
}

// Ladder represents a single ladder
type Ladder struct {
	Name     string         `json:"name"`
	ID       string         `json:"id"`
	Created  time.Time      `json:"created"`
	OwnerKey *datastore.Key `json:"ownerKey"`
	Players  []LadderPlayer `datastore:"noindex" json:"players"`
}

// PlayerLadders represents the ladders a player either owns or is playing in
type PlayerLadders struct {
	Owned   []*Ladder `json:"owned"`
	Playing []*Ladder `json:"playing"`
}

// NewLadder creates a new ladder
func NewLadder() *Ladder {
	return &Ladder{
		Created: time.Now(),
		ID:      xid.New().String(),
		Players: make([]LadderPlayer, 0),
	}
}

const initialRating = 1000

// GetLadder gets a ladder by ID
func GetLadder(ctx context.Context, id string) (*Ladder, error) {
	l := &Ladder{ID: id}
	key := l.DatastoreKey(ctx)
	err := datastore.Get(ctx, key, l)

	if err != nil {
		return nil, err
	}

	return l, nil
}

func GetLaddersForPlayer(ctx context.Context, token *auth.Token) (*PlayerLadders, error) {
	owned := make([]*Ladder, 0)
	playing := make([]*Ladder, 0)

	key := PlayerKeyFromToken(ctx, token)

	_, err := datastore.NewQuery(LadderKind).Filter("OwnerKey = ", key).GetAll(ctx, &owned)

	if err != nil {
		log.Errorf(ctx, "error querying owned ladders for %v: %v", key, err)
		return nil, err
	}

	_, err = datastore.NewQuery(LadderKind).Filter("Players.Key = ", key).GetAll(ctx, &playing)

	if err != nil {
		log.Errorf(ctx, "error querying playing ladders for %v: %v", key, err)
		return nil, err
	}

	// Set empty players to empty array, ensures `[]` rather than `null` in JSON response
	for _, ladder := range append(owned, playing...) {
		if ladder.Players == nil {
			ladder.Players = make([]LadderPlayer, 0)
		}
	}

	return &PlayerLadders{
		Owned:   owned,
		Playing: playing,
	}, nil
}

func (l *Ladder) SetOwner(ctx context.Context, firebaseID string) error {
	p, err := GetPlayer(ctx, firebaseID)

	if err != nil {
		return err
	}

	l.OwnerKey = p.DatastoreKey(ctx)

	return l.AddPlayer(ctx, p)
}

func (l *Ladder) ContainsPlayer(ctx context.Context, p *Player) bool {
	for _, q := range l.Players {
		if p.DatastoreKey(ctx).String() == q.Key.String() {
			return true
		}
	}

	return false
}

func (l *Ladder) AddPlayer(ctx context.Context, p *Player) error {
	key := p.DatastoreKey(ctx)
	_, err := GetPlayer(ctx, p.FirebaseID)

	if err != nil {
		return fmt.Errorf("Error loading Player %s from DB: %s. Will not add to ladder", p.FirebaseID, err.Error())
	}

	if l.ContainsPlayer(ctx, p) {
		return fmt.Errorf("Ladder %s (ID: %s) already contains player with ID %s", l.Name, l.ID, p.FirebaseID)
	}

	lp := LadderPlayer{
		Key:      key,
		Position: len(l.Players) + 1,
		Name:     p.Name,
		Wins:     0,
		Losses:   0,
		Rating:   p.Rating,
	}

	l.Players = append(l.Players, lp)

	return nil
}

// LogGame registers a game in this ladder:
// - writes the Game to the DB as a descendant of this Ladder
// - updates boh Player's wins/losses in this ladder
// - if the winner was previously ranked below the loser, swaps the player positions
// - writes the ladder to the DB with the new results
func (l *Ladder) LogGame(ctx context.Context, g *Game) (*Game, error) {
	var winner, loser *LadderPlayer // winner and loser indexes
	winnerP, loserP := g.WinnerAndLoser()

	for i := 0; i < len(l.Players); i++ {
		p := l.Players[i]

		if p.Key.Equal(winnerP.DatastoreKey(ctx)) {
			winner = &l.Players[i]
		} else if p.Key.Equal(loserP.DatastoreKey(ctx)) {
			loser = &l.Players[i]
		}
	}

	if winner == nil {
		return nil, fmt.Errorf("could not locate game winner %s", winnerP.DatastoreKey(ctx))
	}

	if loser == nil {
		return nil, fmt.Errorf("could not locate game loser %s", loserP.DatastoreKey(ctx))
	}

	// 1. rank the winner and loser and write the updated Players to the DB
	// 2. set the rating change in this Game for both Players
	// 3. save the Game
	// 4. update the ladder statistics for both LadderPlayers and write the
	//    updated ladder to the DB
	err := datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		wa, la := rank(ctx, &winnerP, &loserP)

		if _, err := winnerP.Save(ctx); err != nil {
			return err
		}

		if _, err := loserP.Save(ctx); err != nil {
			return err
		}

		if err := g.SetRatingChange(winnerP, wa); err != nil {
			return err
		}

		if err := g.SetRatingChange(loserP, la); err != nil {
			return err
		}

		if err := g.Save(ctx, l); err != nil {
			return err
		}

		winner.Wins = winner.Wins + 1
		loser.Losses = loser.Losses + 1

		winner.Rating = wa
		loser.Rating = la

		// swap positions if the winner was positioned lower (greater number) than the loser
		if winner.Position > loser.Position {
			winner.Position, loser.Position = loser.Position, winner.Position
		}

		if _, err := l.Save(ctx); err != nil {
			return err
		}

		return nil
	}, &datastore.TransactionOptions{XG: true})

	if err != nil {
		return nil, err
	}

	return g, nil
}

func (l *Ladder) sortPlayers() {
	sort.Slice(l.Players, func(i, j int) bool {
		return l.Players[i].Position < l.Players[j].Position
	})
}

// Games registered against this ladder
func (l *Ladder) Games(ctx context.Context) ([]*Game, error) {
	games := make([]*Game, 0)
	query := datastore.NewQuery(GameKind).Ancestor(l.DatastoreKey(ctx))
	_, err := query.GetAll(ctx, &games)

	if err != nil {
		return nil, err
	}

	return games, nil
}

func (l *Ladder) DatastoreKey(ctx context.Context) *datastore.Key {
	return datastore.NewKey(ctx, LadderKind, l.ID, 0, nil)
}

// Save the ladder to the DB.
// Players are always sorted by their position before saving.
func (l *Ladder) Save(ctx context.Context) (*datastore.Key, error) {
	if l.Name == "" {
		return nil, errors.New("Cannot save a Ladder without a Name")
	}

	key := l.DatastoreKey(ctx)
	l.sortPlayers()
	return datastore.Put(ctx, key, l)
}
