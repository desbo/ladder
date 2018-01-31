package main

import (
	"fmt"
	"log"
	"time"

	firebase "firebase.google.com/go"
	"github.com/desbo/ladder/api"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"

	"google.golang.org/api/option"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/remote_api"
)

// backfill player.SignupDate from firebase
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

	ps := make([]*app.Player, 0)
	q := datastore.NewQuery("Player")

	_, err = q.GetAll(remoteCtx, &ps)

	if err == nil {
		app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile("../api/firebase_key.json"))

		if err != nil {
			log.Fatal(err)
		}

		auth, err := app.Auth(remoteCtx)

		if err != nil {
			log.Fatal(err)
		}

		for _, p := range ps {
			u, err := auth.GetUser(remoteCtx, p.FirebaseID)

			if err != nil {
				log.Fatal(err)
			}

			ts := u.UserMetadata.CreationTimestamp
			d := time.Unix(ts/1000, ts%1000)
			p.SignupDate = d

			_, err = p.Save(remoteCtx)

			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("saved", p.Name)
		}
	}
}
