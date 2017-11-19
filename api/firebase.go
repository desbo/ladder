package app

import (
	firebase "firebase.google.com/go"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

func initFirebase() (*firebase.App, error) {
	opt := option.WithCredentialsFile("firebase_key.json")
	return firebase.NewApp(context.Background(), nil, opt)
}
