package ladder

import (
	"time"

	"google.golang.org/appengine/datastore"
)

type playerScore struct {
	Player *datastore.Key
	Score  int
}

type Game struct {
	Date    time.Time
	Ladder  *datastore.Key
	Player1 playerScore
	Player2 playerScore
}
