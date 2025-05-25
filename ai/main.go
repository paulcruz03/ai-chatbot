package ai

import (
	"context"
	"log"

	"go-chatbot/utils"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func Main() *genai.GenerativeModel {
	aiKey := utils.GoDotEnvVariable("GEMINI_API_KEY")
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(aiKey))

	if err != nil {
		log.Fatal(err)
	}
	model := client.GenerativeModel("gemini-2.0-flash-exp")

	CreateLogger("[AI] Init Ai")
	return model
}

func AiChatPrompt(model *genai.GenerativeModel, chatHistory []*genai.Content, prompt string) string {
	ctx := context.Background()
	cs := model.StartChat()

	CreateLogger("[AI] Prompt: " + prompt)
	cs.History = chatHistory
	resp, err := cs.SendMessage(ctx, genai.Text(prompt))
	if err != nil {
		log.Fatal(err)
	}
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		if textPart, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
			return string(textPart)
		}
	}
	return ""
}

func AiPrompt(model *genai.GenerativeModel, prompt string) string {
	ctx := context.Background()
	finalPrompt := GeneratePrompt(prompt)
	resp, err := model.GenerateContent(ctx, genai.Text(finalPrompt))

	if err != nil {
		log.Fatal(err)
	}
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		if textPart, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
			return string(textPart)
		}
	}
	return ""
}
