package app

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/xid"
	"google.golang.org/appengine/datastore"
)

type playerResult struct {
	Player       Player  `json:"player"`
	Score        int     `json:"score"`
	RatingChange float64 `json:"ratingChange"`
}

type Game struct {
	ID      string       `json:"id"`
	Date    time.Time    `json:"submitted"`
	Player1 playerResult `json:"player1"`
	Player2 playerResult `json:"player2"`
}

func newPlayerResult(p *Player, score int) playerResult {
	return playerResult{
		Player: *p,
		Score:  score,
	}
}

// WinnerAndLoser returns the player Keys for the winner and loser of this match, respectively
func (g *Game) WinnerAndLoser() (Player, Player) {
	if g.Player1.Score > g.Player2.Score {
		return g.Player1.Player, g.Player2.Player
	}

	return g.Player2.Player, g.Player1.Player
}

func (g *Game) Save(ctx context.Context, ladder *Ladder) error {
	key := datastore.NewKey(ctx, GameKind, g.ID, 0, ladder.DatastoreKey(ctx))

	if _, err := datastore.Put(ctx, key, g); err != nil {
		return fmt.Errorf("error saving game %s: %s", g.ID, err)
	}

	return nil
}

// SavePlayers saves both players in this game to the DB
func (g *Game) SavePlayers(ctx context.Context) error {
	winner, loser := g.WinnerAndLoser()

	if _, err := winner.Save(ctx); err != nil {
		return fmt.Errorf("error saving game winner (id %s): %s", winner.FirebaseID, err.Error())
	} else if _, err := loser.Save(ctx); err != nil {
		return fmt.Errorf("error saving game loser (id %s): %s", loser.FirebaseID, err.Error())
	}

	return nil
}

func (g *Game) SetRatingChange(p Player, change float64) error {
	if g.Player1.Player.Equals(p) {
		g.Player1.RatingChange = change
	} else if g.Player2.Player.Equals(p) {
		g.Player2.RatingChange = change
	} else {
		return fmt.Errorf("player %s not in game %s", p.FirebaseID, g.ID)
	}

	return nil
}

// NewGame creates a new game
func NewGame(p1 *Player, p2 *Player, p1score int, p2score int) *Game {
	return &Game{
		ID:      xid.New().String(),
		Date:    time.Now(),
		Player1: newPlayerResult(p1, p1score),
		Player2: newPlayerResult(p2, p2score),
	}
}
