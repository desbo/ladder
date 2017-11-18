package app

import (
	"encoding/json"
	"net/http"

	firebase "firebase.google.com/go"
	"github.com/desbo/ladder"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/appengine"
)

func login(firebase *firebase.App) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		d := new(playerPostData)
		err := decode(d, r)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := appengine.NewContext(r)
		p, err := ladder.GetPlayer(ctx, d.Name)

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if p.CheckPassword(ctx, d.RawPassword) {
			client, err := firebase.Auth(ctx)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			token, err := client.CustomToken(p.Name)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			json.NewEncoder(w).Encode(token)
		} else {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
		}
	}
}
