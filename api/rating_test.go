package app

import (
	"fmt"
	"testing"

	"google.golang.org/appengine/aetest"
)

func mkPlayer(n string) *Player {
	return &Player{
		FirebaseID:       n,
		Name:             n,
		Rating:           1000,
		RatingDeviation:  350,
		RatingVolatility: 0.06,
	}
}

func TestRating(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer done()

	n := 100

	p1 := mkPlayer("timo boll")
	p2 := mkPlayer("ma long")

	for i := 0; i < n; i++ {
		g := NewGame(ctx, p1, p2, 11, 5)
		p1.Rating, p2.Rating, _ = rank(ctx, g)
		fmt.Println(g.Player1.Player.Rating, g.Player2.Player.Rating)
	}
}
