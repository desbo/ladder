package api

import (
	"time"
)

type Point struct {
	Time   time.Time `json:"x"`
	Rating int       `json:"y"`
}

type Line []Point

// A Chart is a map of usernames to a series of Points
type Chart map[string]Line

func BuildChart(games []*Game) Chart {
	chart := make(Chart)

	updateChart := func(time time.Time, username string, ratingChange int) {
		var previous Point

		line, ok := chart[username]

		if !ok {
			line = make(Line, 0)
			chart[username] = line
			previous = Point{Rating: initialRating}
		} else {
			previous = chart[username][len(chart[username])-1]
		}

		chart[username] = append(chart[username], Point{
			Time:   time,
			Rating: previous.Rating + ratingChange,
		})
	}

	for _, game := range games {
		updateChart(game.Date, game.Player1.User.Name, game.Player1.RatingChange)
		updateChart(game.Date, game.Player2.User.Name, game.Player2.RatingChange)
	}

	return chart
}
