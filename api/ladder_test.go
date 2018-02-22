package api

import (
	"fmt"
	"testing"

	"golang.org/x/net/context"
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"
)

const LadderKey string = "test"
const LadderSize = 2

var TestUser *User

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
	t.Run("Get ladders for player", func(t *testing.T) { GetLaddersTest(ctx, t) })
}

func CreateLadderTest(ctx context.Context, t *testing.T) {
	owner := &User{
		FirebaseID: "owner",
		Name:       "owner",
	}

	_, err := owner.Save(ctx)

	if err != nil {
		t.Fatalf("error saving ladder owner %v: %s", owner, err.Error())
	}

	l, err := NewLadder(ctx, owner)

	if err != nil {
		t.Fatalf("error creating ladder: %s", err.Error())
	}

	l.ID = LadderKey
	l.Name = "test ladder"
	l.Save(ctx)

	key := datastore.NewKey(ctx, LadderKind, l.ID, 0, nil)
	result := &Ladder{}

	if err := datastore.Get(ctx, key, result); err != nil {
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

		user := &User{
			FirebaseID: fmt.Sprintf("%d", i),
			Name:       name,
		}

		_, err := user.Save(ctx)

		if err != nil {
			t.Fatalf("error saving user %v: %s", user, err.Error())
		}

		err = l.AddUser(ctx, user)

		if err != nil {
			t.Fatalf("error after adding user %v: %s", user, err.Error())
		}

		if !l.ContainsUser(ctx, user) {
			t.Fatalf("ladder was not reported to contain user %v after adding", user)
		}

		if TestUser == nil {
			TestUser = user
		}
	}

	// +1 because we added owner at the start
	if len(l.Players) != ladderSize+1 {
		t.Fatalf("incorrect numbers of players in ladder: expected %d, got %d", ladderSize, len(l.Players))
	}

	_, err = l.Save(ctx)

	if err != nil {
		t.Fatalf("error saving ladder: %s", err.Error())
	}
}

func GetLaddersTest(ctx context.Context, t *testing.T) {
	owner, err := GetUser(ctx, "owner")

	if err != nil {
		t.Fatalf("could not load ladder owner")
	}

	ls, err := GetLaddersForUser(ctx, TestUser)

	if err != nil {
		t.Fatalf("error getting player ladders: %s", err.Error())
	}

	if len(ls.Playing) != 1 {
		t.Fatalf("incorrect number of playing ladders for %s (expected %d, got %d)", TestUser, 1, len(ls.Playing))
	}

	ls, err = GetLaddersForUser(ctx, owner)

	if len(ls.Owned) != 1 {
		t.Fatalf("incorrect number of playing ladders for %s (expected %d, got %d)", owner, 1, len(ls.Owned))
	}

}

func SubmitGameTest(ctx context.Context, t *testing.T) {
	l, err := GetLadder(ctx, LadderKey)

	if err != nil {
		t.Fatalf("could not get ladder with ID %s: %s", LadderKey, err.Error())
	}

	users := make([]*User, len(l.Players))

	for i := 0; i < 2; i++ {
		users[i], err = UserFromPlayer(ctx, l.Players[i])

		if err != nil {
			t.Fatalf("could not create Player from LadderPlayer %v: %s", l.Players[i], err.Error())
		}
	}

	// match should result in a swap
	winner := users[1]
	loser := users[0]
	game := NewGame(winner, loser, 11, 5)

	_, err = l.LogGame(ctx, game)

	if err != nil {
		t.Fatalf("error logging game: %s", err.Error())
	}

	games, err := l.GamesForCurrentSeason(ctx)

	if err != nil {
		t.Fatalf("error looking up games for ladder: %s", err.Error())
	}

	if len(games) != 1 {
		t.Fatalf("wrong number of games in ladder, got %d, expected 1", len(games))
	}

	if l.Players[0].Name != winner.Name {
		t.Fatalf("position of winner set incorrectly")
	}

	if l.Players[len(l.Players)-1].Name != loser.Name {
		t.Fatalf("position of loser set incorrectly")
	}

	if len(l.Players) != LadderSize+1 {
		t.Fatalf("wrong number of players in ladder after game: got %d, expected %d", len(l.Players), LadderSize)
	}

	winner, err = GetUserByEncodedKey(ctx, winner.DatastoreKey(ctx).Encode())

	if err != nil {
		t.Fatal(err)
	}
}
