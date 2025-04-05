package handlers

import (
	"net/http"

	"go-chatbot/ai"

	"github.com/gin-gonic/gin"
)

func GetAiResponse(c *gin.Context) {
	client := ai.Main()
	if client == nil {
		c.Abort()
	}

	c.JSON(http.StatusOK, gin.H{
		"response": ai.AiPrompt(client, "What is Ai?"),
	})
}
