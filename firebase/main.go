package backend

import (
	"context"
	"fmt"
	"go-chatbot/utils"
	"log"
	"time"

	"go-chatbot/ai"

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

type ChatHistory struct {
	Role string `json:"role,omitempty"`
	Text string `json:"text,omitempty"`
}

type Chat struct {
	CreatedAt int64          `json:"createdAt,omitempty"`
	UpdatedAt int64          `json:"updatedAt,omitempty"`
	History   []*ChatHistory `json:"history,omitempty"`
	Title     string         `json:"title,omitempty"`
}

func New() (*Backend, error) {
	firebaseServiceAccountJSON := utils.GoDotEnvVariable("FIREBASE_SERVICE_ACCOUNT_JSON")
	realtimeDbUrl := utils.GoDotEnvVariable("FIREBASE_REALTIME_DB")
	if firebaseServiceAccountJSON == "" {
		log.Fatalf("FIREBASE_SERVICE_ACCOUNT_JSON environment variable not set")
	}
	if realtimeDbUrl == "" {
		log.Fatalf("FIREBASE_REALTIME_DB environment variable not set")
	}
	conf := &firebase.Config{
		DatabaseURL: realtimeDbUrl,
	}

	opt := option.WithCredentialsJSON([]byte(firebaseServiceAccountJSON))
	app, err := firebase.NewApp(context.Background(), conf, opt)
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
		return nil, false
	}
	log.Printf("Verified ID token for UID: %s\n", token.UID)
	return token, true
}

func (b Backend) CreateChat(userUID string, userPrompt string) (*db.Ref, *Chat, error) {
	log.Printf("Creating new chat for %s", userUID)
	ctx := context.Background()
	ai, err := ai.New("")
	if err != nil {
		log.Fatalln("Error Ai Init:", err)
		return nil, nil, err
	}

	ref := b.dbClient.NewRef(fmt.Sprintf(`%s/chats`, userUID))

	newPostRef, err := ref.Push(ctx, nil)
	if err != nil {
		log.Fatalln("Error pushing child node:", err)
		return nil, nil, err
	}

	history := []*ChatHistory{}
	history = append(history, &ChatHistory{Role: genai.RoleUser, Text: userPrompt})
	history = append(history, &ChatHistory{Role: genai.RoleModel, Text: ai.GenerateSingleResponse(userPrompt)})
	chat := Chat{
		CreatedAt: time.Now().UTC().UnixMilli(),
		Title:     ai.GenerateTitle(userPrompt),
		History:   history,
	}

	if err := newPostRef.Set(ctx, chat); err != nil {
		log.Fatalln("Error setting value:", err)
	}

	return newPostRef, &chat, nil
}

func (b Backend) VerifyAndRetrieveChat(userUID string, key string) (*Chat, []*genai.Content, bool) {
	log.Printf("Getting chat history for %s", userUID)
	ctx := context.Background()
	ref := b.dbClient.NewRef(fmt.Sprintf(`%s/chats/%s`, userUID, key))

	var chatHistory *Chat
	if err := ref.Get(ctx, &chatHistory); err != nil {
		log.Printf("Error getting post %s: %v\n", key, err)
		return nil, nil, false
	}

	if chatHistory == nil {
		return nil, nil, false
	}

	convertedHistory := []*genai.Content{}
	for _, chat := range chatHistory.History {
		convertedHistory = append(convertedHistory, genai.NewContentFromText(string(chat.Text), genai.Role(chat.Role)))
	}
	return chatHistory, convertedHistory, true
}

func (b Backend) UpdateChat(userUID string, key string, newHistory []*genai.Content) {
	history := []*ChatHistory{}

	for _, chat := range newHistory {
		history = append(history, &ChatHistory{
			Role: chat.Role,
			Text: chat.Parts[0].Text,
		})
	}

	ctx := context.Background()
	ref := b.dbClient.NewRef(fmt.Sprintf(`%s/chats/%s`, userUID, key))
	historyRef := ref.Child("history")
	err := historyRef.Set(ctx, history)

	if err != nil {
		log.Fatalln("Error setting value:", err)
	}
}
