package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"google.golang.org/appengine"
)

type playerPostData struct {
	Name        string `json:"name"`
	RawPassword string `json:"password"`
}

func getLadder(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	lad, err := GetLadder(appengine.NewContext(r), ps.ByName("id"))

	if err != nil {
		fmt.Fprintf(w, "%s", err)
	} else {
		json.NewEncoder(w).Encode(lad)
	}
}

func getLaddersForPlayer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := appengine.NewContext(r)
	token, err := initAndVerifyToken(ctx, r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ladders, err := GetLaddersForPlayer(ctx, token)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ladders)
}

func createLadder(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ladder := NewLadder()
	err := decode(ladder, r)
	ctx := appengine.NewContext(r)

	if ladder.Name == "" {
		http.Error(w, "no ladder name provided", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := initAndVerifyToken(ctx, r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ladder.OwnerKey = PlayerKeyFromToken(ctx, token)

	save(ladder, w, r)
}

func createPlayer(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := appengine.NewContext(r)
	token, err := initAndVerifyToken(ctx, r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	form := struct {
		Name string `json:"name"`
	}{}

	err = json.NewDecoder(r.Body).Decode(&form)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	save(NewPlayer(token, form.Name), w, r)
}

func init() {
	router := httprouter.New()

	router.GET("/ladder/:id", getLadder)
	router.POST("/ladder", createLadder)

	router.GET("/ladders", getLaddersForPlayer)

	router.POST("/player", createPlayer)

	http.Handle("/", cors.AllowAll().Handler(router))
}
