package handlers

import (
	"fmt"
	"go-chatbot/ai"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var AiModel = ai.Main()
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

		resp := ai.AiChatPrompt(AiModel, chatSessions[clientId], string(msg))
		ws.WriteMessage(t, []byte(resp))

		// create new chat history
		chatSessions[clientId] = append(chatSessions[clientId], &genai.Content{
			Parts: []genai.Part{
				genai.Text(string(msg)),
			},
			Role: "user",
		}, &genai.Content{
			Parts: []genai.Part{
				genai.Text(resp),
			},
			Role: "model",
		})
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
