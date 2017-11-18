package ladder

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// Player is a user in the system that can create and join ladders
type Player struct {
	Name     string `json:"name"`
	Password string `json:"-"`
	Results  []struct {
		LadderName string
		Wins       int
		Losses     int
		LadderKey  *datastore.Key
	} `datastore:",noindex" json:"results"`
}

// NewPlayer creates a player with the name and encrypted password
func NewPlayer(name string, rawPassword string) (*Player, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	return &Player{
		Name:     name,
		Password: string(password[:]),
	}, nil
}

// Save a player to the DB
func (p *Player) Save(ctx context.Context) (*datastore.Key, error) {
	key := datastore.NewKey(ctx, "Player", p.Name, 0, nil)

	err := datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		err := datastore.Get(ctx, key, p)

		if err == nil {
			return fmt.Errorf("player %s already exists", p.Name)
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
