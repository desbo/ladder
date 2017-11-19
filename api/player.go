package app

import (
	"errors"
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
	if rawPassword == "" {
		return nil, errors.New("no password provided")
	}

	password, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	return &Player{
		Name:     name,
		Password: string(password[:]),
	}, nil
}

func key(ctx context.Context, id string) *datastore.Key {
	return datastore.NewKey(ctx, "Player", id, 0, nil)
}

func (p *Player) Key(ctx context.Context) *datastore.Key {
	return key(ctx, p.Name)
}

func (p *Player) CheckPassword(ctx context.Context, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.Password), []byte(password))

	if err != nil {
		return false
	}

	return true
}

func GetPlayer(ctx context.Context, name string) (*Player, error) {
	p := new(Player)
	err := datastore.Get(ctx, key(ctx, name), p)

	if err != nil {
		return nil, err
	}

	return p, nil
}

// Save a player to the DB
func (p *Player) Save(ctx context.Context) (*datastore.Key, error) {
	key := p.Key(ctx)

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
