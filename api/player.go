package app

import (
	"fmt"
	"net/http"

	"firebase.google.com/go/auth"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// Player is a user in the system that can create and join ladders
// FirebaseID is used as the datastore key
type Player struct {
	FirebaseID       string  `json:"-"`
	Name             string  `json:"name"`
	Rating           float64 `json:"rating"`
	RatingDeviation  float64 `json:"-"`
	RatingVolatility float64 `json:"-"`
}

// NewPlayer creates a new player
// Initial rating values are based on http://www.glicko.net/glicko/glicko2.pdf
func NewPlayer(token *auth.Token, name string) *Player {
	return &Player{
		FirebaseID:       token.UID,
		Name:             name,
		Rating:           1000,
		RatingDeviation:  350,
		RatingVolatility: 0.06,
	}
}

func GetPlayer(ctx context.Context, firebaseID string) (*Player, error) {
	p := &Player{}
	err := datastore.Get(ctx, datastore.NewKey(ctx, PlayerKind, firebaseID, 0, nil), p)

	if err != nil {
		return nil, err
	}

	return p, nil
}

func GetPlayerByEncodedKey(ctx context.Context, key string) (*Player, error) {
	p := &Player{}
	dec, err := datastore.DecodeKey(key)

	if err != nil {
		return nil, err
	}

	if err = datastore.Get(ctx, dec, p); err != nil {
		return nil, err
	}

	return p, nil
}

func PlayerKeyFromToken(ctx context.Context, token *auth.Token) *datastore.Key {
	return datastore.NewKey(ctx, PlayerKind, token.UID, 0, nil)
}

func PlayerFromLadderPlayer(ctx context.Context, lp LadderPlayer) (*Player, error) {
	p := &Player{}
	err := datastore.Get(ctx, lp.Key, p)

	if err != nil {
		return nil, err
	}

	return p, nil
}

func PlayerFromRequest(r *http.Request) (*Player, error) {
	ctx := appengine.NewContext(r)
	token, err := initAndVerifyToken(ctx, r)

	if err != nil {
		return nil, err
	}

	return GetPlayer(ctx, token.UID)
}

func (p *Player) DatastoreKey(ctx context.Context) *datastore.Key {
	return datastore.NewKey(ctx, PlayerKind, p.FirebaseID, 0, nil)
}

func (p *Player) Equals(o Player) bool {
	return p.FirebaseID == o.FirebaseID
}

// Save a player to the DB
func (p *Player) Save(ctx context.Context) (*datastore.Key, error) {
	key := p.DatastoreKey(ctx)

	err := datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		err := datastore.Get(ctx, key, p)

		if err == nil {
			return fmt.Errorf("player %s already exists", p.FirebaseID)
		} else if err != datastore.ErrNoSuchEntity {
			return err
		}

		_, err = datastore.Put(ctx, key, p)

		return err
	}, nil)

	if err != nil {
		return nil, err
	}

	return key, nil
}
