package signaling

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type WebsocketConnectionManager struct {
	connections map[string]*websocket.Conn
}

func (manager *WebsocketConnectionManager) Add(identifier string, websocket *websocket.Conn) {
	if _, exist := manager.connections[identifier]; exist {
		fmt.Printf("%v is already in the connection pool", identifier)
		return
	}
	manager.connections[identifier] = websocket
}

func (manager *WebsocketConnectionManager) Delete(identifier string) {
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
