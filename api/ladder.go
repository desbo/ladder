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
	Rating   float64        `json:"rating"`
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

// GetLadder gets a ladder from an encoded Datastore key
func GetLadder(ctx context.Context, encodedKey string) (*Ladder, error) {
	l := &Ladder{}

	key, err := datastore.DecodeKey(encodedKey)

	if err != nil {
		return nil, err
	}

	err = datastore.Get(ctx, key, l)

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
func (l *Ladder) LogGame(ctx context.Context, g *Game) error {

	var winner, loser LadderPlayer
	winnerP, loserP := g.WinnerAndLoser()

	// retrieve winner and loser from LadderPlayers
	// remove them both now, to be added back after they're updated
	for i := 0; i < len(l.Players); i++ {
		p := l.Players[i]
		remove := false

		if p.Key.Equal(winnerP.DatastoreKey(ctx)) {
			winner = p
			remove = true
		} else if p.Key.Equal(loserP.DatastoreKey(ctx)) {
			loser = p
			remove = true
		}

		if remove {
			err := l.removeNthPlayer(i)

			if err != nil {
				return err
			}
		}
	}

	wr, lr, err := rank(ctx, g)

	if err != nil {
		return err
	}

	key := datastore.NewKey(ctx, GameKind, g.ID, 0, l.DatastoreKey(ctx))

	if _, err = datastore.Put(ctx, key, g); err != nil {
		return err
	}

	winner.Wins = winner.Wins + 1
	winner.Rating = wr
	loser.Losses = loser.Losses + 1
	loser.Rating = lr

	// swap positions if the winner was ranked lower (greater number) than the loser
	if winner.Position > loser.Position {
		winner.Position, loser.Position = loser.Position, winner.Position
	}

	l.Players = append(l.Players, winner, loser)
	_, err = l.Save(ctx)

	return err
}

func (l *Ladder) removeNthPlayer(n int) error {
	ps := l.Players

	if n >= len(ps) {
		return fmt.Errorf("index out of bounds when removing player: index %d, len(players) %d", n, len(ps))
	}

	// see https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-array-in-golang/37335777
	ps[len(ps)-1], ps[n] = ps[n], ps[len(ps)-1]
	l.Players = ps[:len(ps)-1]

	return nil
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
