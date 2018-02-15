package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	// enable remote API
	_ "google.golang.org/appengine/remote_api"
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

	user, err := GetUserFromToken(ctx, token)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ladders, err := GetLaddersForUser(ctx, user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ladders)
}

func createLadder(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := appengine.NewContext(r)
	user, err := UserFromRequest(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ladder, err := NewLadder(ctx, user)

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

func createUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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

	save(NewUser(token, form.Name), w, r)
}

func joinLadder(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := appengine.NewContext(r)
	user, err := UserFromRequest(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	lad, err := GetLadder(ctx, p.ByName("id"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = lad.AddUser(ctx, user); err != nil {
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
	user, err := UserFromRequest(r)

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

	opponent, err := GetUserByEncodedKey(ctx, form.OpponentKey)

	if err != nil {
		log.Errorf(ctx, "error getting ladder %s: %s", form.LadderID, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	game := NewGame(user, opponent, form.MyScore, form.TheirScore)
	ladder, err := GetLadder(ctx, form.LadderID)

	if err != nil {
		log.Errorf(ctx, "error logging game %s (could not get ladder %s: %s)", game.ID, form.LadderID, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if game, err = ladder.LogGame(ctx, game); err != nil {
		log.Errorf(ctx, "error logging game %s: %s", game.ID, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Infof(ctx, "logged game %s", game.ID)

	json.NewEncoder(w).Encode(game)
}

func chart(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := appengine.NewContext(r)
	l, err := GetLadder(ctx, p.ByName("id"))

	if err != nil {
		log.Errorf(ctx, "error generating chart (could not get ladder %s: %s)", p.ByName("id"), err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	games, err := l.Games(ctx)

	if err != nil {
		log.Errorf(ctx, "error generating chart (could not get games for ladder %s: %s)", p.ByName("id"), err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(BuildChart(games))
}

func inactive(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if err := CheckInactivity(appengine.NewContext(r)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode("OK")
}

func jsonify(f func(w http.ResponseWriter, r *http.Request, p httprouter.Params)) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		f(w, r, p)
	}
}

func init() {
	router := httprouter.New()

	router.POST("/game", jsonify(submitGame))

	router.POST("/join/:id", jsonify(joinLadder))

	router.GET("/ladder/:id", jsonify(getLadder))
	router.POST("/ladder", jsonify(createLadder))
	router.GET("/ladders", jsonify(getLaddersForPlayer))

	router.POST("/player", jsonify(createUser))

	router.GET("/chart/:id", jsonify(chart))

	router.GET("/cron/inactive", inactive)

	http.Handle("/", cors.AllowAll().Handler(router))
}
