package api

import (
	"net/http"
	"time"

	"firebase.google.com/go/auth"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// User is a user in the system that can create and join ladders
// FirebaseID is used as the datastore key
type User struct {
	FirebaseID string    `json:"-"`
	Name       string    `json:"name"`
	SignupDate time.Time `json:"signupDate"`
}

// NewUser creates a new User
// Initial rating values are based on http://www.glicko.net/glicko/glicko2.pdf
func NewUser(token *auth.Token, name string) *User {
	return &User{
		FirebaseID: token.UID,
		Name:       name,
		SignupDate: time.Now(),
	}
}

func GetUser(ctx context.Context, firebaseID string) (*User, error) {
	p := &User{}
	err := datastore.Get(ctx, datastore.NewKey(ctx, UserKind, firebaseID, 0, nil), p)

	if err != nil {
		return nil, err
	}

	return p, nil
}

func GetUserByEncodedKey(ctx context.Context, key string) (*User, error) {
	p := &User{}
	dec, err := datastore.DecodeKey(key)

	if err != nil {
		return nil, err
	}

	if err = datastore.Get(ctx, dec, p); err != nil {
		return nil, err
	}

	return p, nil
}

func GetUserFromToken(ctx context.Context, token *auth.Token) (*User, error) {
	p := &User{}
	key := datastore.NewKey(ctx, UserKind, token.UID, 0, nil)

	if err := datastore.Get(ctx, key, p); err != nil {
		return nil, err
	}

	return p, nil
}

func UserFromPlayer(ctx context.Context, p Player) (*User, error) {
	u := &User{}
	err := datastore.Get(ctx, p.Key, u)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func UserFromRequest(r *http.Request) (*User, error) {
	ctx := appengine.NewContext(r)
	token, err := initAndVerifyToken(ctx, r)

	if err != nil {
		return nil, err
	}

	return GetUser(ctx, token.UID)
}

func (u *User) DatastoreKey(ctx context.Context) *datastore.Key {
	return datastore.NewKey(ctx, UserKind, u.FirebaseID, 0, nil)
}

func (u *User) Equals(o User) bool {
	return u.FirebaseID == o.FirebaseID
}

func (u *User) Games(ctx context.Context) ([]*Game, error) {
	games := make([]*Game, 0)

	if _, err := datastore.NewQuery(GameKind).
		Filter("Player1.User.FirebaseID = ", u.FirebaseID).
		GetAll(ctx, &games); err != nil {
		return nil, err
	}

	if _, err := datastore.NewQuery(GameKind).
		Filter("Player2.User.FirebaseID = ", u.FirebaseID).
		GetAll(ctx, &games); err != nil {
		return nil, err
	}

	return games, nil
}

func (u *User) LastActive(ctx context.Context) (*time.Time, error) {
	games, err := u.Games(ctx)

	if err != nil {
		return nil, err
	}

	t := u.SignupDate

	for _, game := range games {
		if game.Date.After(t) {
			t = game.Date
		}
	}

	return &t, nil
}

// Save a player to the DB
func (u *User) Save(ctx context.Context) (*datastore.Key, error) {
	key := u.DatastoreKey(ctx)

	if _, err := datastore.Put(ctx, key, u); err != nil {
		return nil, err
	}

	return key, nil
}
