package handlers

import (
	"fmt"
	"go-chatbot/ai"
	backend "go-chatbot/firebase"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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

) {
	defer ws.Close()

	chatHistory := chat.History
	for {
		t, msg, err := ws.ReadMessage()
		if err != nil {
			break
		}

		client, err := ai.New(chatId)
		if err != nil {
			break
		}
		resp, history := client.Send(string(msg), chatHistory)
		ws.WriteMessage(t, []byte(resp))

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

	chatHistory, doesChatExist := firebase.VerifyAndRetrieveChat(c.Param("clientId"), c.Param("chatId"))
	if !doesChatExist {
		return
	}

	ws, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	chat(ws, firebase, c.Param("clientId"), c.Param("chatId"), chatHistory)
}
