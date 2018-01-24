package app

import (
	"fmt"
	"net/http"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"golang.org/x/net/context"
)

const firebaseTokenHeader = "Firebase-Token"

func initFirebase(ctx context.Context) (*firebase.App, error) {
	return firebase.NewApp(ctx, nil)
}

func getFirebaseToken(r *http.Request) string {
	return r.Header.Get(firebaseTokenHeader)
}

func verifyToken(ctx context.Context, app *firebase.App, tokenStr string) (*auth.Token, error) {
	client, err := app.Auth(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting firebase.Auth client: %v", err)
	}

	return client.VerifyIDToken(tokenStr)
}

func initAndVerifyToken(ctx context.Context, r *http.Request) (*auth.Token, error) {
	firebase, err := initFirebase(ctx)

	if err != nil {
		return nil, err
	}

	return verifyToken(ctx, firebase, getFirebaseToken(r))
}
