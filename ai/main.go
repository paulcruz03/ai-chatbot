package ai

import (
	"context"
	"fmt"
	"go-chatbot/utils"
	"log"

	"google.golang.org/genai"
)

type AiClient struct {
	ID     string
	Client *genai.Client
	Model  string
}

func New(id string) (*AiClient, error) {
	aiKey := utils.GoDotEnvVariable("GEMINI_API_KEY")

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  aiKey,
		Backend: genai.BackendGeminiAPI,
	})

	if err != nil {
		return nil, err
	}

	model := "gemini-2.0-flash-exp" // or any other model you want to use

	return &AiClient{
		ID:     id,
		Client: client,
		Model:  model,
	}, nil
}

func (e AiClient) Send(msg string, history []*genai.Content) (string, []*genai.Content) {
	ctx := context.Background()

	chat, _ := e.Client.Chats.Create(ctx, e.Model, nil, history)
	res, _ := chat.SendMessage(ctx, genai.Part{Text: msg})
	fmt.Println(e.ID, chat)

	history = append(history,
		genai.NewContentFromText(msg, genai.RoleUser),
		genai.NewContentFromText(res.Candidates[0].Content.Parts[0].Text, genai.RoleModel),
	)

	if len(res.Candidates) > 0 {
		return res.Candidates[0].Content.Parts[0].Text, history
	}

	return "", history
}

func (e AiClient) GenerateTitle(msg string) string {
	ctx := context.Background()
	result, err := e.Client.Models.GenerateContent(
		ctx,
		e.Model,
		genai.Text(fmt.Sprintf(`
			Generate a concise, engaging, and relevant title for the following user prompt: %s

			Only return one title
		`, msg)),
		nil,
	)
	if err != nil {
		log.Fatalf("error generating title: %v\n", err)
		return msg
	}

	return result.Text()
}
