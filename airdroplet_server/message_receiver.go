package signaling

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type MessageHandler struct {
	conn              *websocket.Conn
	connectionManager *WebsocketConnectionManager
}

func (handler *MessageHandler) handle() {
	for {
		_, message, err := handler.conn.ReadMessage()
		if err != nil {
			return
		}
		fmt.Println(string(message))
	}
}
