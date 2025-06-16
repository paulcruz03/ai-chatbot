package main

import (
	"go-chatbot/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	router := gin.Default()
	router.Use(CORSMiddleware())

	router.LoadHTMLFiles("ws_tester.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "ws_tester.html", gin.H{})
	})

	router.GET("/health", handlers.HealthCheck)
	// router.GET("/chat-list", )
	// router.GET("/chat-init")

	router.GET("/ai", handlers.GetAiResponse)
	router.GET("/client-id", handlers.GenerateClientId)
	router.GET("/ws/:clientId", handlers.WsHandler)

	router.Run(":8080")
	// Start HTTPS server
	// err := router.RunTLS(":80", "./cert.pem", "./key.pem")
	// if err != nil {
	// 	panic("Failed to start HTTPS server: " + err.Error())
	// }
}
