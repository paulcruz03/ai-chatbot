package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go-chatbot/ai"
	"go-chatbot/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
)

func GetAiResponse(c *gin.Context) {
	client := ai.Main()
	if client == nil {
		c.Abort()
	}

	c.JSON(http.StatusOK, gin.H{
		"response": ai.AiChatPrompt(client, []*genai.Content{}, "What is Ai?"),
	})
}

func GetAiQuestions(c *gin.Context) {
	client := ai.Main()
	if client == nil {
		c.Abort()
	}

	mode := c.Query("mode")

	if mode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Prompt and output format are required"})
		return
	}

	prompt := ""
	outputFormat := ""

	switch mode {
	case "game-init-questions":
		{
			promptObj, err := utils.GoGetJsonValue("prompt.json", "workout-init-plan")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing AI response"})
				return
			}
			prompt = promptObj.Prompt
			outputFormat = promptObj.ExpectedOutput.(string)
		}
	default:
		{
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing AI response"})
			return
		}
	}

	rawResponse := ai.AiPrompt(client, prompt, outputFormat)
	var response map[string]interface{}
	err := json.Unmarshal([]byte(rawResponse), &response)
	if err != nil {
		fmt.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing AI response"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"response": response})
}
