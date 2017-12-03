package app

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// Entity kinds
const (
	LadderKind = "Ladder"
	PlayerKind = "Player"
)

// DatastoreEntity is any entity that can exist in GAE datastore
type DatastoreEntity interface {
	Save(ctx context.Context) (*datastore.Key, error)
}

type saveResult struct {
	DatastoreEntity `json:"entity"`
	Key             *datastore.Key `json:"key"`
}

func decode(e interface{}, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	return decoder.Decode(e)
}

func save(e DatastoreEntity, w http.ResponseWriter, r *http.Request) {
	key, err := e.Save(appengine.NewContext(r))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(saveResult{
		DatastoreEntity: e,
		Key:             key,
	})
}
