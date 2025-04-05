package ai

import (
	"context"
	"log"

	"go-chatbot/utils"

	"google.golang.org/genai"
)

func Main() *genai.Client {
	aiKey := utils.GoDotEnvVariable("GEMINI_API_KEY")
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{APIKey: aiKey})

	if err != nil {
		log.Fatal(err)
	}

	CreateLogger("Init Ai") // helper function for printing content parts
	return client
}

func AiPrompt(client *genai.Client, prompt string) string {
	ctx := context.Background()
	parts := []*genai.Part{
		{Text: prompt},
	}
	resp, err := client.Models.GenerateContent(ctx, "gemini-2.0-flash-exp", []*genai.Content{{Parts: parts}}, nil)
	if err != nil {
		log.Fatal(err)
	}
	CreateLogger("Ai Prompt: " + prompt) // helper function for printing content parts
	return resp.Text()
}
