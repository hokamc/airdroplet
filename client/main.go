package main

import (
	"net"
	"signaling"
)

func main() {
	server := signaling.NewWebServer(net.IPv4(127, 0, 0, 1), 8080)
	websocketConnection := signaling.NewWebsocketConnection()
	server.AddRoute("/ws", websocketConnection.HandleWebsocketConnection)
	server.Start()
}
