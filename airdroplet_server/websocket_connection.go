package signaling

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

const ReadBufferSize = 1024
const WriteBufferSize = 1024

type WebsocketConnection struct {
	connectionManager *WebsocketConnectionManager
	upgrader          websocket.Upgrader
}

func NewWebsocketConnection() *WebsocketConnection {
	websocketConnection := new(WebsocketConnection)
	websocketConnection.upgrader = websocket.Upgrader{
		ReadBufferSize:  ReadBufferSize,
		WriteBufferSize: WriteBufferSize,
	}
	websocketConnection.connectionManager = NewWebsocketConnectionManger()
	return websocketConnection
}

func (server *WebsocketConnection) HandleWebsocketConnection(response http.ResponseWriter, request *http.Request) {
	id := request.URL.Query()["id"]
	if id == nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	conn, _ := server.upgrader.Upgrade(response, request, nil)
	server.connectionManager.add(id[0], conn)

	handler := MessageHandler{conn: conn, connectionManager: server.connectionManager}
	go handler.handle(request.Context())
}

type WebsocketConnectionManager struct {
	connections map[string]*websocket.Conn
}

func (manager *WebsocketConnectionManager) add(identifier string, websocket *websocket.Conn) {
	if _, exist := manager.connections[identifier]; exist {
		fmt.Printf("%v is already in the connection pool", identifier)
		return
	}
	manager.connections[identifier] = websocket
}

func (manager *WebsocketConnectionManager) delete(identifier string) {
	if _, exist := manager.connections[identifier]; exist {
		manager.connections[identifier].Close()
		delete(manager.connections, identifier)
	}
}

func NewWebsocketConnectionManger() *WebsocketConnectionManager {
	manager := new(WebsocketConnectionManager)
	manager.connections = make(map[string]*websocket.Conn)
	return manager
}
