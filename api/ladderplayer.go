package app

// LadderPlayer is a player in a ladder
import (
	"context"
	"sort"
	"time"

	"google.golang.org/appengine/datastore"
)

type LadderPlayer struct {
	Key      *datastore.Key `json:"key" datastore:"__key__"` // firebase ID
	Name     string         `json:"name"`
	Position int            `json:"position"`
	Wins     int            `json:"wins"`
	Losses   int            `json:"losses"`
	Rating   int            `json:"rating"`
	WinRate  float32        `json:"winRate"`
	JoinDate time.Time      `json:"joinDate"`
}

type LadderPlayers []LadderPlayer

func NewLadderPlayer(ctx context.Context, p *Player, position int) LadderPlayer {
	return LadderPlayer{
		Key:      p.DatastoreKey(ctx),
		Position: position,
		Name:     p.Name,
		Wins:     0,
		Losses:   0,
		Rating:   p.Rating,
		JoinDate: time.Now(),
	}
}

func (lp *LadderPlayer) winRate() float32 {
	if lp.Wins == 0 {
		return 0
	}

	return float32(lp.Wins) / float32(lp.Wins+lp.Losses)
}

// implement sort.Interface to sort players by rating
func (lps LadderPlayers) Len() int {
	return len(lps)
}

func (lps LadderPlayers) Less(i, j int) bool {
	return lps[i].Rating < lps[j].Rating
}

func (lps LadderPlayers) Swap(i, j int) {
	lps[i], lps[j] = lps[j], lps[i]
}

func (lps LadderPlayers) sortByRanking() {
	sort.Sort(sort.Reverse(lps))

	for i := 0; i < len(lps); i++ {
		lps[i].Position = i + 1
	}
}
