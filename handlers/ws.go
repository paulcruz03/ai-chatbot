package handlers

import (
	"fmt"
	"go-chatbot/ai"
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
var chatSessions map[string][]*genai.Content

func chat(clientId string, ws *websocket.Conn) {
	defer ws.Close()

	// to create a new chat session
	// like { clientId: chatHistory[] }
	if chatSessions == nil {
		chatSessions = make(map[string][]*genai.Content)
	}

	for {
		t, msg, err := ws.ReadMessage()
		if err != nil {
			break
		}

		client, err := ai.New(clientId)
		if err != nil {
			break
		}
		resp, history := client.Send(string(msg), chatSessions[clientId])
		ws.WriteMessage(t, []byte(resp))

		// create new chat history
		chatSessions[clientId] = history
	}
}

func WsHandler(ctx *gin.Context) {
	// security check that only allowed clientId can access the websocket
	if !CheckAllowedClientId(ctx.Param("clientId")) {
		return
	}

	ws, err := wsupgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	chat(ctx.Param("clientId"), ws)
}
