package handlers

import (
	"go-chatbot/ai"
	backend "go-chatbot/firebase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserPrompt struct {
	IdToken         string `json:"idToken"`
	UserPromptValue string `json:"userPrompt"`
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "up",
	})
}

func StartChat(c *gin.Context) {
	var userPrompt UserPrompt
	if err := c.BindJSON(&userPrompt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	firebase, err := backend.New()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, isVerified := firebase.Verify(userPrompt.IdToken)
	if !isVerified {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	ai, err := ai.New("")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	chatId, chatContent, err := firebase.CreateChat(token.UID, ai.GenerateTitle(userPrompt.UserPromptValue))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"chatId": chatId.Key,
		"title":  chatContent.Title,
	})
}
