package handlers

import (
	"fmt"
	"go-chatbot/ai"
	"net/http"
	"time"

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

var AiClient = ai.Main()

func chat(ws *websocket.Conn) {
	defer ws.Close()
	for {
		t, msg, err := ws.ReadMessage()
		if err != nil {
			break
		}
		resp := ai.AiPrompt(AiClient, string(msg))
		time.Sleep(5000 * time.Millisecond)
		ws.WriteMessage(t, []byte(resp))
	}
}

func WsHandler(ctx *gin.Context) {
	if !CheckAllowedClientId(ctx.Param("clientId")) {
		return
	}

	ws, err := wsupgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	chat(ws)
}
