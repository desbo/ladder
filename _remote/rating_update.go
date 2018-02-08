package main

import (
	"fmt"
	"log"

	"github.com/desbo/ladder/api"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/remote_api"
)

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

	ps := make([]*api.Player, 0)
	ls := make([]*api.Ladder, 0)

	lq := datastore.NewQuery("Ladder")
	pq := datastore.NewQuery("Player")

	_, err = lq.GetAll(remoteCtx, &ls)
	_, err = pq.GetAll(remoteCtx, &ps)

	if err == nil {
		for _, l := range ls {
			for i, _ := range l.Players {
				l.Players[i].Rating = l.Players[i].Rating + 500
			}

			_, err := l.Save(remoteCtx)

			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("saved", l)
		}
	}

	for _, p := range ps {
		p.Rating = p.Rating + 500

		_, err = p.Save(remoteCtx)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("saved", p.Name)
	}
}
