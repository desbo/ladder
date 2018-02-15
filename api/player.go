package api

// Player is a player in a ladder
import (
	"context"
	"sort"
	"time"

	"google.golang.org/appengine/datastore"
)

// Player is a player in a ladder. The same User may have different ratings in
// different ladders.
type Player struct {
	Key              *datastore.Key `json:"key" datastore:"__key__"` // firebase ID
	Name             string         `json:"name"`
	Position         int            `json:"position"`
	Wins             int            `json:"wins"`
	Losses           int            `json:"losses"`
	Rating           int            `json:"rating"`
	RatingDeviation  float64        `json:"-"`
	RatingVolatility float64        `json:"-"`
	WinRate          float32        `json:"winRate"`
	JoinDate         time.Time      `json:"joinDate"`
	Active           bool           `json:"active"`
}

type Players []Player

func NewPlayer(ctx context.Context, u *User, position int) Player {
	return Player{
		Key:              u.DatastoreKey(ctx),
		Position:         position,
		Name:             u.Name,
		Wins:             0,
		Losses:           0,
		Rating:           initialRating,
		RatingDeviation:  350,
		RatingVolatility: 0.06,
		JoinDate:         time.Now(),
		Active:           true,
	}
}

func (lp *Player) winRate() float32 {
	if lp.Wins == 0 {
		return 0
	}

	return float32(lp.Wins) / float32(lp.Wins+lp.Losses)
}

// implement sort.Interface to sort players by rating
func (ps Players) Len() int {
	return len(ps)
}

func (ps Players) Less(i, j int) bool {
	return ps[i].Rating < ps[j].Rating
}

func (ps Players) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

func (ps Players) sortByRanking() {
	sort.Sort(sort.Reverse(ps))

	pos := 1

	for i := 0; i < len(ps); i++ {
		if !ps[i].Active {
			ps[i].Position = 0
		} else {
			ps[i].Position = pos
			pos++
		}
	}
}
