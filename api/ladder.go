package app

import (
	"errors"
	"fmt"
	"time"

	"firebase.google.com/go/auth"

	"github.com/gosimple/slug"
	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

const maxIDIncrements = 5

type LadderPlayer struct {
	Key      *datastore.Key `json:"key" datastore:"__key__"`
	Position int            `json:"position"`
	Name     string         `json:"name"`
	Wins     int            `json:"wins"`
	Losses   int            `json:"losses"`
}

// Ladder represents a single ladder
type Ladder struct {
	ID       string         `json:"id"`
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

	return l, nil
}

func GetLaddersForPlayer(ctx context.Context, token *auth.Token) (*PlayerLadders, error) {
	owned := make([]*Ladder, 0)
	playing := make([]*Ladder, 0)

	key := PlayerKeyFromToken(ctx, token)

	_, err := datastore.NewQuery(LadderKind).Filter("OwnerKey = ", key).GetAll(ctx, &owned)

	if err != nil {
		log.Errorf(ctx, "error querying owned ladders for %v: %v", key, err)
		return nil, err
	}

	_, err = datastore.NewQuery(LadderKind).Filter("Players.Key = ", key).GetAll(ctx, &playing)

	if err != nil {
		log.Errorf(ctx, "error querying playing ladders for %v: %v", key, err)
		return nil, err
	}

	return &PlayerLadders{
		Owned:   owned,
		Playing: playing,
	}, nil
}

func generateID(ctx context.Context, l *Ladder, suffix int) string {
	id := slug.Make(l.Name)

	if suffix > 0 {
		id = fmt.Sprintf("%s-%d", id, suffix)
	}

	return id
}

func attemptSave(ctx context.Context, l *Ladder, attempt int) func(ctx context.Context) error {
	if attempt >= maxIDIncrements {
		return func(ctx context.Context) error {
			return errors.New("Unable to save ladder (ID conflict)")
		}
	}

	id := generateID(ctx, l, attempt)
	key := datastore.NewKey(ctx, LadderKind, id, 0, nil)
	err := datastore.Get(ctx, key, &Ladder{})

	if err == nil {
		return attemptSave(ctx, l, attempt+1)
	} else if err != datastore.ErrNoSuchEntity {
		return func(ctx context.Context) error {
			return err
		}
	}

	l.ID = id

	_, err = datastore.Put(ctx, key, l)

	return func(ctx context.Context) error {
		return err
	}
}

// Save the ladder to the DB
func (l *Ladder) Save(ctx context.Context) (*datastore.Key, error) {
	if l.Name == "" {
		return nil, errors.New("Attempted to save a Ladder that has no name")
	}

	if l.Players == nil {
		l.Players = make([]LadderPlayer, 0)
	}

	err := datastore.RunInTransaction(ctx, attemptSave(ctx, l, 0), nil)

	if err != nil {
		return nil, err
	}

	return datastore.NewKey(ctx, LadderKind, l.ID, 0, nil), nil
}
