package handlers

import (
	"encoding/json"
	"fmt"
	"go-chatbot/ai"
	backend "go-chatbot/firebase"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"google.golang.org/genai"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func chat(
	ws *websocket.Conn,
	backend *backend.Backend,
	userUid string,
	chatId string,
	chat *backend.Chat,
	chatHistory []*genai.Content,
) {
	defer ws.Close()

	for {
		t, msg, err := ws.ReadMessage()
		if err != nil {
			break
		}
		loadingMsg := `{"type":"loading","message":"AI is typing..."}`
		ws.WriteMessage(t, []byte(loadingMsg))

		client, err := ai.New(chatId)
		if err != nil {
			log.Println("Error:", err)
			break
		}
		resp, history := client.Send(string(msg), chatHistory)
		responseMsg, err := json.Marshal(map[string]string{"type": "response", "message": resp})
		if err != nil {
			log.Println("JSON marshal error:", err)
			continue // or break
		}
		ws.WriteMessage(t, responseMsg)

		// create new chat history
		chatHistory = history
	}

	log.Printf("Closing uid: %s , chatId: %s\n", userUid, chatId)
	backend.UpdateChat(userUid, chatId, chatHistory)
}

func WsHandler(c *gin.Context) {
	firebase, err := backend.New()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	chatHistory, convertedHistory, doesChatExist := firebase.VerifyAndRetrieveChat(
		c.Param("clientId"),
		c.Param("chatId"),
	)
	if !doesChatExist {
		return
	}

	ws, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	chat(ws, firebase, c.Param("clientId"), c.Param("chatId"), chatHistory, convertedHistory)
}
