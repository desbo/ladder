package app

import (
	"fmt"

	"firebase.google.com/go/auth"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// Player is a user in the system that can create and join ladders
// FirebaseID is used as the datastore key
type Player struct {
	FirebaseID string `datastore:",noindex" json:"-"`
	Results    []struct {
		LadderName string
		Wins       int
		Losses     int
		LadderKey  *datastore.Key
	} `datastore:",noindex" json:"results"`
}

func NewPlayerFromToken(token *auth.Token) *Player {
	return &Player{
		FirebaseID: token.UID,
	}
}

func PlayerKeyFromToken(ctx context.Context, token *auth.Token) *datastore.Key {
	return datastore.NewKey(ctx, PlayerKind, token.UID, 0, nil)
}

// Save a player to the DB
func (p *Player) Save(ctx context.Context) (*datastore.Key, error) {
	key := datastore.NewKey(ctx, PlayerKind, p.FirebaseID, 0, nil)

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
