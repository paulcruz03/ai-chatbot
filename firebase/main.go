package backend

import (
	"context"
	"go-chatbot/utils"
	"log"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"firebase.google.com/go/v4/db"
	"google.golang.org/api/option"
	"google.golang.org/genai"
)

type Backend struct {
	client     *firebase.App
	authClient *auth.Client
	dbClient   *db.Client
}

type Chat struct {
	CreatedAt int64            `json:"createdAt,omitempty"`
	History   []*genai.Content `json:"history,omitempty"`
	Title     string           `json:"title,omitempty"`
}

func New() (*Backend, error) {
	firebaseServiceAccountJSON := utils.GoDotEnvVariable("FIREBASE_SERVICE_ACCOUNT_JSON")

	if firebaseServiceAccountJSON == "" {
		log.Fatalf("FIREBASE_SERVICE_ACCOUNT_JSON environment variable not set")
	}
	// fmt.Print(firebaseServiceAccountJSON)

	opt := option.WithCredentialsJSON([]byte(firebaseServiceAccountJSON))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
		return nil, err
	}

	authClient, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
		return nil, err
	}

	dbClient, err := app.Database(context.Background())
	if err != nil {
		log.Fatalln("Error initializing database client:", err)
		return nil, err
	}

	return &Backend{
		client:     app,
		authClient: authClient,
		dbClient:   dbClient,
	}, nil
}

func (b Backend) Verify(idToken string) (*auth.Token, bool) {
	token, err := b.authClient.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		log.Fatalf("error verifying ID token: %v\n", err)
		return token, true
	}
	log.Printf("Verified ID token for UID: %s\n", token.UID)
	return nil, false
}

func (b Backend) CreateChat(userUID string, title string) (*db.Ref, *Chat, error) {
	ctx := context.Background()
	ref := b.dbClient.NewRef(userUID)
	postsRef := ref.Child("chats")

	newPostRef, err := postsRef.Push(ctx, nil)
	if err != nil {
		log.Fatalln("Error pushing child node:", err)
		return nil, nil, err
	}

	chat := Chat{
		CreatedAt: time.Now().UTC().UnixMilli(),
		Title:     title,
	}
	if err := newPostRef.Set(ctx, chat); err != nil {
		log.Fatalln("Error setting value:", err)
	}

	return newPostRef, &chat, nil
}
