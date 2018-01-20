package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

type playerPostData struct {
	Name        string `json:"name"`
	RawPassword string `json:"password"`
}

func getLadder(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	lad, err := GetLadder(appengine.NewContext(r), p.ByName("id"))

	if err != nil {
		fmt.Fprintf(w, "error getting ladder: %s", err)
	} else {
		json.NewEncoder(w).Encode(lad)
	}
}

func getLaddersForPlayer(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := appengine.NewContext(r)
	token, err := initAndVerifyToken(ctx, r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	player, err := GetPlayerFromToken(ctx, token)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ladders, err := GetLaddersForPlayer(ctx, player)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ladders)
}

func createLadder(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := appengine.NewContext(r)
	player, err := PlayerFromRequest(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ladder, err := NewLadder(ctx, player)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := decode(ladder, r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !ladder.Valid(ctx) {
		// TODO: improve error reporting
		http.Error(w, "Invalid ladder", http.StatusBadRequest)
		return
	}

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

func joinLadder(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := appengine.NewContext(r)
	player, err := PlayerFromRequest(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	lad, err := GetLadder(ctx, p.ByName("id"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = lad.AddPlayer(ctx, player); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := lad.Save(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode("OK")
}

func submitGame(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := appengine.NewContext(r)
	userPlayer, err := PlayerFromRequest(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	form := struct {
		LadderID    string `json:"ladderID"`
		MyScore     int    `json:"myScore"`
		TheirScore  int    `json:"theirScore"`
		OpponentKey string `json:"opponentKey"`
	}{}

	defer r.Body.Close()

	if err = json.NewDecoder(r.Body).Decode(&form); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Errorf(ctx, "error getting opponent %s: %s", form.OpponentKey, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	opponent, err := GetPlayerByEncodedKey(ctx, form.OpponentKey)

	if err != nil {
		log.Errorf(ctx, "error getting ladder %s: %s", form.LadderID, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	game := NewGame(userPlayer, opponent, form.MyScore, form.TheirScore)
	ladder, err := GetLadder(ctx, form.LadderID)

	if err != nil {
		log.Errorf(ctx, "error getting ladder %s: %s", form.LadderID, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if game, err = ladder.LogGame(ctx, game); err != nil {
		log.Errorf(ctx, "error logging game %s: %s", game, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(game)
}

func init() {
	router := httprouter.New()

	router.POST("/game", submitGame)
	router.POST("/join/:id", joinLadder)
	router.GET("/ladder/:id", getLadder)
	router.POST("/ladder", createLadder)
	router.GET("/ladders", getLaddersForPlayer)
	router.POST("/player", createPlayer)

	http.Handle("/", cors.AllowAll().Handler(router))
}
