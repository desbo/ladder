package app

import (
	"time"

	"firebase.google.com/go/auth"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

type LadderPlayer struct {
	Key      *datastore.Key `json:"key" datastore:"__key__"`
	Position int            `json:"position"`
	Name     string         `json:"name"`
	Wins     int            `json:"wins"`
	Losses   int            `json:"losses"`
}

// Ladder represents a single ladder
type Ladder struct {
	Key      *datastore.Key `json:"key"`
	Name     string         `json:"name"`
	Created  time.Time      `json:"created"`
	OwnerKey *datastore.Key `json:"ownerKey"`
	Players  []LadderPlayer `datastore:",noindex" json:"players"`
}

// PlayerLadders represents the ladders a player either owns or is playing in
type PlayerLadders struct {
	Owned   []*Ladder `json:"owned"`
	Playing []*Ladder `json:"playing"`
}

// NewLadder creates a new ladder
func NewLadder() *Ladder {
	return &Ladder{
		Created: time.Now(),
		Players: make([]LadderPlayer, 0),
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

	l.Key = key

	return l, nil
}

func GetLaddersForPlayer(ctx context.Context, token *auth.Token) (*PlayerLadders, error) {
	owned := make([]*Ladder, 0)
	playing := make([]*Ladder, 0)

	key := PlayerKeyFromToken(ctx, token)

	ownedKeys, err := datastore.NewQuery(LadderKind).Filter("OwnerKey = ", key).GetAll(ctx, &owned)

	if err != nil {
		log.Errorf(ctx, "error querying owned ladders for %v: %v", key, err)
		return nil, err
	}

	playerKeys, err := datastore.NewQuery(LadderKind).Filter("Players.Key = ", key).GetAll(ctx, &playing)

	if err != nil {
		log.Errorf(ctx, "error querying playing ladders for %v: %v", key, err)
		return nil, err
	}

	keys := append(ownedKeys, playerKeys...)

	// Set empty players to empty array, ensures `[]` rather than `null` in JSON response
	for i, ladder := range append(owned, playing...) {
		ladder.Key = keys[i]

		if ladder.Players == nil {
			ladder.Players = make([]LadderPlayer, 0)
		}
	}

	return &PlayerLadders{
		Owned:   owned,
		Playing: playing,
	}, nil
}

// Save the ladder to the DB
func (l *Ladder) Save(ctx context.Context) (*datastore.Key, error) {
	key := datastore.NewIncompleteKey(ctx, LadderKind, nil)
	return datastore.Put(ctx, key, l)
}
