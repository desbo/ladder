package ladder

import (
	"time"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
)

// Ladder represents a single ladder
type Ladder struct {
	Name    string         `json:"name"`
	Created time.Time      `json:"created"`
	Owner   *datastore.Key `json:"owner"`
	Players []struct {
		Position int
		Name     string
		Wins     int
		Losses   int
		Key      *datastore.Key
	} `datastore:",noindex" json:"players"`
}

// NewLadder creates a new ladder
func NewLadder() *Ladder {
	return &Ladder{
		Created: time.Now(),
	}
}

// GetLadder gets a ladder
func GetLadder(ctx context.Context, id string) (*Ladder, error) {
	l := &Ladder{}

	key, err := datastore.DecodeKey(id)

	if err != nil {
		return nil, err
	}

	err = datastore.Get(ctx, key, l)

	if err != nil {
		return nil, err
	}

	return l, nil
}

// Save the ladder to the DB
func (l *Ladder) Save(ctx context.Context) (*datastore.Key, error) {
	key := datastore.NewIncompleteKey(ctx, "Ladder", nil)
	return datastore.Put(ctx, key, l)
}
