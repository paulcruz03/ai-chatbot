package handlers

import (
	"net/http"

	"go-chatbot/ai"

	"github.com/gin-gonic/gin"
	"google.golang.org/genai"
)

func GetAiResponse(c *gin.Context) {
	client, err := ai.New("test-client")
	if err != nil {
		c.Abort()
	}

	message, _ := client.Send("Hello, how are you?", []*genai.Content{})
	c.JSON(http.StatusOK, gin.H{
		"response": message,
	})
}
