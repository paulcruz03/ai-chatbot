package handlers

import (
	"crypto/rand"
	"net/http"

	"go-chatbot/ai"

	"github.com/gin-gonic/gin"
	"google.golang.org/genai"
)

var allowedClientIds = []string{""}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "up",
	})
}

func GenerateClientId(c *gin.Context) {
	id := rand.Text()
	allowedClientIds = append(allowedClientIds, id)

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

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
