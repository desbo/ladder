package app

import (
	"net/http"
	"time"

	"firebase.google.com/go/auth"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// Player is a user in the system that can create and join ladders
// FirebaseID is used as the datastore key
type Player struct {
	FirebaseID       string    `json:"-"`
	Name             string    `json:"name"`
	Rating           int       `json:"rating"`
	RatingDeviation  float64   `json:"-"`
	RatingVolatility float64   `json:"-"`
	SignupDate       time.Time `json:"signupDate"`
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
		SignupDate:       time.Now(),
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

func GetPlayerFromToken(ctx context.Context, token *auth.Token) (*Player, error) {
	p := &Player{}
	key := datastore.NewKey(ctx, PlayerKind, token.UID, 0, nil)

	if err := datastore.Get(ctx, key, p); err != nil {
		return nil, err
	}

	return p, nil
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

func (p *Player) Games(ctx context.Context) ([]*Game, error) {
	games := make([]*Game, 0)

	if _, err := datastore.NewQuery(GameKind).
		Filter("Player1.Player.FirebaseID = ", p.FirebaseID).
		GetAll(ctx, games); err != nil {
		return nil, err
	}

	if _, err := datastore.NewQuery(GameKind).
		Filter("Player2.Player.FirebaseID = ", p.FirebaseID).
		GetAll(ctx, games); err != nil {
		return nil, err
	}

	return games, nil
}

func (p *Player) LastActive(ctx context.Context) (*time.Time, error) {
	games, err := p.Games(ctx)

	if err != nil {
		return nil, err
	}

	t := p.SignupDate

	for _, game := range games {
		if game.Date.After(t) {
			t = game.Date
		}
	}

	return &t, nil
}

// Save a player to the DB
func (p *Player) Save(ctx context.Context) (*datastore.Key, error) {
	key := p.DatastoreKey(ctx)

	if _, err := datastore.Put(ctx, key, p); err != nil {
		return nil, err
	}

	return key, nil
}
