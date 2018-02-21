package api

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/net/context"

	"github.com/rs/xid"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

// Ladder represents a single ladder
type Ladder struct {
	Name             string         `json:"name"`
	ID               string         `json:"id"`
	Created          time.Time      `json:"created"`
	OwnerKey         *datastore.Key `json:"ownerKey"`
	Players          Players        `json:"players"`
	InactivityPeriod time.Duration  `json:"-"`
	Active           bool           `json:"active"`
	Season           int            `json:"season"`
}

// LaddersForUser represents the ladders a user either owns or is playing in
type LaddersForUser struct {
	Owned   []*Ladder `json:"owned"`
	Playing []*Ladder `json:"playing"`
}

const initialRating = 1500

// NewLadder creates a new ladder
func NewLadder(ctx context.Context, owner *User) (*Ladder, error) {
	l := &Ladder{
		ID:               xid.New().String(),
		Created:          time.Now(),
		Active:           true,
		Players:          make(Players, 0),
		OwnerKey:         owner.DatastoreKey(ctx),
		InactivityPeriod: 7 * 24 * time.Hour, // TODO: Add to UI,
		Season:           1,
	}

	if err := l.AddUser(ctx, owner); err != nil {
		return nil, err
	}

	return l, nil
}

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

func GetLaddersForUser(ctx context.Context, user *User) (*LaddersForUser, error) {
	owned := make([]*Ladder, 0)
	playing := make([]*Ladder, 0)

	key := user.DatastoreKey(ctx)

	_, err := datastore.NewQuery(LadderKind).Filter("OwnerKey = ", key).GetAll(ctx, &owned)

	if err != nil {
		log.Errorf(ctx, "error querying owned ladders for %v: %v", key, err)
		return nil, err
	}

	_, err = datastore.NewQuery(LadderKind).Filter("Players.Name = ", user.Name).GetAll(ctx, &playing)

	if err != nil {
		log.Errorf(ctx, "error querying playing ladders for %v: %v", key, err)
		return nil, err
	}

	// Set empty players to empty array, ensures `[]` rather than `null` in JSON response
	for _, ladder := range append(owned, playing...) {
		if ladder.Players == nil {
			ladder.Players = make(Players, 0)
		}
	}

	return &LaddersForUser{
		Owned:   owned,
		Playing: playing,
	}, nil
}

func (l *Ladder) ContainsUser(ctx context.Context, u *User) bool {
	for _, q := range l.Players {
		if u.DatastoreKey(ctx).String() == q.Key.String() {
			return true
		}
	}

	return false
}

func (l *Ladder) AddUser(ctx context.Context, u *User) error {
	_, err := GetUser(ctx, u.FirebaseID)

	if err != nil {
		return fmt.Errorf("Error loading Player %s from DB: %s. Will not add to ladder", u.FirebaseID, err.Error())
	}

	if l.ContainsUser(ctx, u) {
		return fmt.Errorf("Ladder %s (ID: %s) already contains player with ID %s", l.Name, l.ID, u.FirebaseID)
	}

	lp := NewPlayer(ctx, u, len(l.Players)+1)
	l.Players = append(l.Players, lp)

	return nil
}

// LogGame registers a game in this ladder:
// - writes the Game to the DB as a descendant of this Ladder
// - updates boh Player's wins/losses in this ladder
// - if the winner was previously ranked below the loser, swaps the player positions
// - writes the ladder to the DB with the new results
func (l *Ladder) LogGame(ctx context.Context, g *Game) (*Game, error) {
	var winner, loser *Player // winner and loser indexes
	winnerUser, loserUser := g.WinnerAndLoser()

	for i := 0; i < len(l.Players); i++ {
		p := l.Players[i]

		if p.Key.Equal(winnerUser.DatastoreKey(ctx)) {
			winner = &l.Players[i]
		} else if p.Key.Equal(loserUser.DatastoreKey(ctx)) {
			loser = &l.Players[i]
		}
	}

	if winner == nil {
		return nil, fmt.Errorf("could not locate game winner %s", winnerUser.DatastoreKey(ctx))
	}

	if loser == nil {
		return nil, fmt.Errorf("could not locate game loser %s", loserUser.DatastoreKey(ctx))
	}

	// 1. rank the winner and loser and write the updated Players to the DB
	// 2. set the rating change in this Game for both Players
	// 3. save the Game
	// 4. update the ladder statistics for both LadderPlayers and write the
	//    updated ladder to the DB
	err := datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		wa, la := rank(ctx, winner, loser)

		if err := g.SetRatingChange(winnerUser, wa); err != nil {
			return err
		}

		if err := g.SetRatingChange(loserUser, la); err != nil {
			return err
		}

		if err := g.Save(ctx, l); err != nil {
			return err
		}

		winner.Active = true
		loser.Active = true
		winner.Wins = winner.Wins + 1
		loser.Losses = loser.Losses + 1
		winner.LastRatingChange = wa
		loser.LastRatingChange = la

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

func (l *Ladder) updateWinRates() {
	for i := 0; i < len(l.Players); i++ {
		l.Players[i].WinRate = l.Players[i].winRate()
	}
}

func (l *Ladder) DatastoreKey(ctx context.Context) *datastore.Key {
	return datastore.NewKey(ctx, LadderKind, l.ID, 0, nil)
}

func (l *Ladder) Valid(ctx context.Context) bool {
	if l.Name == "" {
		log.Errorf(ctx, "ladder %s had no Name", l)
		return false
	}

	if l.ID == "" {
		log.Errorf(ctx, "ladder %s had no ID", l)
		return false
	}

	if l.OwnerKey == nil {
		log.Errorf(ctx, "ladder %s had no OwnerKey", l)
		return false
	}

	return true
}

// Save the ladder to the DB.
// Players are always sorted by their position before saving.
func (l *Ladder) Save(ctx context.Context) (*datastore.Key, error) {
	if l.Name == "" {
		return nil, errors.New("Cannot save a Ladder without a Name")
	}

	key := l.DatastoreKey(ctx)

	l.updateWinRates()
	l.Players.sortByRanking()

	if !l.Valid(ctx) {
		return nil, fmt.Errorf("Invalid Ladder %v", l)
	}

	return datastore.Put(ctx, key, l)
}
