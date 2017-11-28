package app

import (
	"google.golang.org/appengine/datastore"
)

// Player is a user in the system that can create and join ladders
type Player struct {
	Name    string `json:"name"`
	Results []struct {
		LadderName string
		Wins       int
		Losses     int
		LadderKey  *datastore.Key
	} `datastore:",noindex" json:"results"`
}
