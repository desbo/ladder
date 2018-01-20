package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/remote_api"

	ladder "github.com/desbo/ladder/api"
)

type LadderPlayerv2 struct {
	Key      *datastore.Key `json:"key" datastore:"__key__"` // firebase ID
	Name     string         `json:"name"`
	Position int            `json:"position"`
	Wins     int            `json:"wins"`
	Losses   int            `json:"losses"`
	Rating   int            `json:"rating"`
}

// Ladder represents a single ladder
type Ladderv2 struct {
	Name     string           `json:"name"`
	ID       string           `json:"id"`
	Created  time.Time        `json:"created"`
	OwnerKey *datastore.Key   `json:"ownerKey"`
	Players  []LadderPlayerv2 `json:"players"`
}

func (l *Ladderv2) Save(ctx context.Context) (*datastore.Key, error) {
	if l.Name == "" {
		return nil, errors.New("Cannot save a Ladder without a Name")
	}

	key := datastore.NewKey(ctx, ladder.LadderKind, l.ID, 0, nil)

	return datastore.Put(ctx, key, l)
}

// Migrate ladder schema to start indexing LadderPlayers
func main() {
	const host = "api-dot-tt-ladder.appspot.com"

	ctx := context.Background()

	hc, err := google.DefaultClient(ctx,
		"https://www.googleapis.com/auth/appengine.apis",
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/cloud-platform",
	)
	if err != nil {
		log.Fatal(err)
	}

	remoteCtx, err := remote_api.NewRemoteContext(host, hc)

	if err != nil {
		log.Fatal(err)
	}

	ls := make([]*ladder.Ladder, 0)

	_, err = datastore.NewQuery(ladder.LadderKind).GetAll(remoteCtx, &ls)

	for _, l := range ls {

		ps2 := make([]LadderPlayerv2, len(l.Players))

		for _, p := range l.Players {
			p2 := LadderPlayerv2{
				Key:      p.Key,
				Name:     p.Name,
				Position: p.Position,
				Wins:     p.Wins,
				Losses:   p.Losses,
				Rating:   p.Rating,
			}

			ps2 = append(ps2, p2)
		}

		l2 := Ladderv2{
			Name:     l.Name,
			ID:       l.ID,
			Created:  l.Created,
			OwnerKey: l.OwnerKey,
			Players:  ps2,
		}

		k, err := l2.Save(remoteCtx)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(k)

	}
}
