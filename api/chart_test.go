package api

import (
	"testing"
	"time"
)

func TestChart(t *testing.T) {
	games := make([]*Game, 0)
	date := time.Now()
	u1 := User{Name: "Timo Boll"}
	u2 := User{Name: "Ma Long"}
	n := 100

	for i := 0; i < n; i++ {
		game := &Game{
			ID:   string(i),
			Date: date.Add(time.Hour * time.Duration(i)),
			Player1: PlayerResult{
				User:         u1,
				RatingChange: 1,
			},
			Player2: PlayerResult{
				User:         u2,
				RatingChange: -1,
			},
		}
		games = append(games, game)
	}

	chart := BuildChart(games)

	timo := chart[u1.Name]
	ma := chart[u2.Name]

	if len(timo) != n || len(ma) != n {
		t.Fatalf("wrong length")
	}

	if timo[len(timo)-1].Rating != initialRating+n {
		t.Fatalf("incorrect rating (%s): expected %d, got %d", u1.Name, initialRating+n, timo[len(timo)-1].Rating)
	}

	if ma[len(ma)-1].Rating != initialRating-n {
		t.Fatalf("incorrect rating (%s): expected %d, got %d", u2.Name, initialRating-n, ma[len(ma)-1].Rating)
	}
}
