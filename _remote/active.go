package main

import (
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

	//ps := make([]*api.Player, 0)
	ls := make([]*api.Ladder, 0)
	q := datastore.NewQuery("Ladder")

	_, err = q.GetAll(remoteCtx, &ls)

	if err == nil {
		for _, l := range ls {
			l.Active = true

			for i, _ := range l.Players {
				l.Players[i].Active = true
			}

			_, err := l.Save(remoteCtx)

			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
