package api

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/xid"
	"google.golang.org/appengine/datastore"
)

type PlayerResult struct {
	User  User `json:"user"`
	Score int  `json:"score"`

	// RatingChange is how much this game altered the player's rating.
	RatingChange int `json:"ratingChange"`
}

type Game struct {
	ID      string       `json:"id"`
	Date    time.Time    `json:"date"`
	Season  int          `json:"-"`
	Player1 PlayerResult `json:"player1"`
	Player2 PlayerResult `json:"player2"`
}

func NewPlayerResult(u *User, score int) PlayerResult {
	return PlayerResult{
		User:  *u,
		Score: score,
	}
}

// WinnerAndLoser returns the winner and loser of this match, respectively
func (g *Game) WinnerAndLoser() (User, User) {
	if g.Player1.Score > g.Player2.Score {
		return g.Player1.User, g.Player2.User
	}

	return g.Player2.User, g.Player1.User
}

func (g *Game) Save(ctx context.Context, ladder *Ladder) error {
	key := datastore.NewKey(ctx, GameKind, g.ID, 0, ladder.DatastoreKey(ctx))

	if _, err := datastore.Put(ctx, key, g); err != nil {
		return fmt.Errorf("error saving game %s: %s", g.ID, err)
	}

	return nil
}

func (g *Game) SetRatingChange(u User, change int) error {
	if g.Player1.User.Equals(u) {
		g.Player1.RatingChange = change
	} else if g.Player2.User.Equals(u) {
		g.Player2.RatingChange = change
	} else {
		return fmt.Errorf("user %s not in game %s", u.FirebaseID, g.ID)
	}

	return nil
}

func (l *Ladder) GamesForSeason(ctx context.Context, season int) ([]*Game, error) {
	games := make([]*Game, 0)
	query := datastore.NewQuery(GameKind).Ancestor(l.DatastoreKey(ctx)).
		Filter("Season =", season).
		Order("Date")

	if _, err := query.GetAll(ctx, &games); err != nil {
		return nil, err
	}

	return games, nil
}

func (l *Ladder) GamesForCurrentSeason(ctx context.Context) ([]*Game, error) {
	return l.GamesForSeason(ctx, l.Season)
}

// NewGame creates a new game
func NewGame(u1 *User, u2 *User, u1score int, u2score int) *Game {
	return &Game{
		ID:      xid.New().String(),
		Date:    time.Now(),
		Player1: NewPlayerResult(u1, u1score),
		Player2: NewPlayerResult(u2, u2score),
	}
}
