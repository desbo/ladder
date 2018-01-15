package app

import (
	"fmt"
	"testing"

	"golang.org/x/net/context"
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"
)

const LadderKey string = "test"
const LadderSize = 2

// TestLadders represents an end-to-end test of:
// - creating a ladder
// - adding players
// - submitting game results
// - updating ladder standings
//
// The same context is used throughout, so database is
// persisted between tests.
func TestLadders(t *testing.T) {

	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer done()

	t.Run("Create ladder", func(t *testing.T) { CreateLadderTest(ctx, t) })
	t.Run("Add players", func(t *testing.T) { AddPlayersTest(ctx, LadderSize, t) })
	t.Run("Submit game", func(t *testing.T) { SubmitGameTest(ctx, t) })
}

func CreateLadderTest(ctx context.Context, t *testing.T) {
	l := NewLadder()
	l.ID = LadderKey
	l.Name = "test ladder"
	l.Save(ctx)

	key := datastore.NewKey(ctx, LadderKind, l.ID, 0, nil)
	result := NewLadder()

	err := datastore.Get(ctx, key, result)

	if err != nil {
		t.Fatal(err)
	}

	if result.ID != l.ID {
		t.Fatalf("ladder was not saved with correct ID: expected %s, got %s", l.ID, result.ID)
	}
}

func AddPlayersTest(ctx context.Context, ladderSize int, t *testing.T) {
	l, err := GetLadder(ctx, LadderKey)

	if err != nil {
		t.Fatalf("could not get ladder with ID %s: %s", LadderKey, err.Error())
	}

	for i := 0; i < ladderSize; i++ {
		name := fmt.Sprintf("Player %d", i)

		player := &Player{
			FirebaseID:       fmt.Sprintf("%d", i),
			Name:             name,
			Rating:           1000,
			RatingDeviation:  350,
			RatingVolatility: 0.06,
		}

		_, err := player.Save(ctx)

		if err != nil {
			t.Fatalf("error saving player %v: %s", player, err.Error())
		}

		err = l.AddPlayer(ctx, player)

		if err != nil {
			t.Fatalf("error after adding player %v: %s", player, err.Error())
		}

		if !l.ContainsPlayer(ctx, player) {
			t.Fatalf("ladder was not reported to contain player %v after adding", player)
		}
	}

	if len(l.Players) != ladderSize {
		t.Fatalf("incorrect numbers of players in ladder: expected %d, got %d", ladderSize, len(l.Players))
	}

	_, err = l.Save(ctx)

	if err != nil {
		t.Fatalf("error saving ladder: %s", err.Error())
	}
}

func SubmitGameTest(ctx context.Context, t *testing.T) {
	l, err := GetLadder(ctx, LadderKey)

	if err != nil {
		t.Fatalf("could not get ladder with ID %s: %s", LadderKey, err.Error())
	}

	players := make([]*Player, len(l.Players))

	for i := 0; i < 2; i++ {
		players[i], err = PlayerFromLadderPlayer(ctx, l.Players[i])

		if err != nil {
			t.Fatalf("could not create Player from LadderPlayer %v: %s", l.Players[i], err.Error())
		}
	}

	// match should result in a swap
	winner := players[1]
	loser := players[0]
	game := NewGame(winner, loser, 11, 5)
	_, err = l.LogGame(ctx, game)

	if err != nil {
		t.Fatalf("error logging game: %s", err.Error())
	}

	games, err := l.Games(ctx)

	if err != nil {
		t.Fatalf("error looking up games for ladder: %s", err.Error())
	}

	if len(games) != 1 {
		t.Fatalf("wrong number of games in ladder, got %d, expected 1", len(games))
	}

	if l.Players[0].Name != winner.Name {
		t.Fatalf("position of winner set incorrectly")
	}

	if l.Players[1].Name != loser.Name {
		t.Fatalf("position of loser set incorrectly")
	}

	if len(l.Players) != LadderSize {
		t.Fatalf("wrong number of players in ladder after game: got %d, expected %d", len(l.Players), LadderSize)
	}
}
