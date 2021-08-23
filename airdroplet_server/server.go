package signaling

import (
	"github.com/gorilla/websocket"
	"net/http"
)

const ReadBufferSize = 1024
const WriteBufferSize = 1024

type WebsocketServer struct {
	connectionManager *WebsocketConnectionManager
	connectionUpgrade websocket.Upgrader
}

func NewWebsocketServer() *WebsocketServer {
	server := new(WebsocketServer)
	server.connectionUpgrade = websocket.Upgrader{
		ReadBufferSize:  ReadBufferSize,
		WriteBufferSize: WriteBufferSize,
	}
	server.connectionManager = NewWebsocketConnectionManger()
	return server
}

// param --> id
func (server *WebsocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := r.Form.Get("id")
	if len(id) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	conn, _ := server.connectionUpgrade.Upgrade(w, r, nil)
	server.connectionManager.Add(id, conn)

	handler := MessageHandler{conn: conn, connectionManager: server.connectionManager}
	go handler.handle()
}
