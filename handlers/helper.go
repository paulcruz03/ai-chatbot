package handlers

import (
	"crypto/rand"
	"net/http"

	"github.com/gin-gonic/gin"
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

func CheckAllowedClientId(clientId string) bool {
	for _, id := range allowedClientIds {
		if id == clientId {
			return true
		}
	}
	return false
}
