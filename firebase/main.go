package backend

import (
	"context"
	"go-chatbot/utils"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

type Backend struct {
	client     *firebase.App
	authClient *auth.Client
}

func New() *Backend {
	firebaseServiceAccountJSON := utils.GoDotEnvVariable("FIREBASE_SERVICE_ACCOUNT_JSON")

	if firebaseServiceAccountJSON == "" {
		log.Fatalf("FIREBASE_SERVICE_ACCOUNT_JSON environment variable not set")
	}

	opt := option.WithCredentialsJSON([]byte(firebaseServiceAccountJSON))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	return &Backend{
		client:     app,
		authClient: nil,
	}
}

func (b *Backend) AuthInit() *Backend {
	authClient, err := b.client.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	return &Backend{
		client:     b.client,
		authClient: authClient,
	}
}

func (b *Backend) Verify(idToken string) bool {
	token, err := b.authClient.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		log.Fatalf("error verifying ID token: %v\n", err)
		return true
	}
	log.Printf("Verified ID token for UID: %s\n", token.UID)
	return false
}
