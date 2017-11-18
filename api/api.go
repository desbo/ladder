package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/desbo/ladder"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/appengine"
)

func getLadder(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	lad, err := ladder.GetLadder(appengine.NewContext(r), ps.ByName("id"))

	if err != nil {
		fmt.Fprintf(w, "%s", err)
	} else {
		json.NewEncoder(w).Encode(lad)
	}
}

func createLadder(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ladder := ladder.NewLadder()
	err := decode(ladder, r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	save(ladder, w, r)
}

func createPlayer(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	type postData struct {
		Name        string `json:"name"`
		RawPassword string `json:"password"`
	}

	d := new(postData)

	err := decode(d, r)
	player, err := ladder.NewPlayer(d.Name, d.RawPassword)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	save(player, w, r)
}

func init() {
	router := httprouter.New()

	router.GET("/ladder/:id", getLadder)
	router.POST("/ladder", createLadder)

	router.POST("/player", createPlayer)

	http.Handle("/", router)
}
