package signaling

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type WebServer struct {
	ip         net.IP
	port       int
	router     *http.ServeMux
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewWebServer(ipAddress net.IP, port int) *WebServer {
	if ipAddress == nil {
		panic("IP address is not specified")
	}

	if port == 0 {
		panic(fmt.Sprintf("Port is not supported, %v", port))
	}

	server := new(WebServer)
	server.ip = ipAddress
	server.port = port
	server.router = http.NewServeMux()
	server.ctx, server.cancelFunc = context.WithCancel(context.Background())

	return server
}

func (server *WebServer) AddRoute(path string, handler func(writer http.ResponseWriter, request *http.Request)) *WebServer {
	validatePath(path)
	server.router.Handle(path, http.HandlerFunc(handler))
	return server
}

// Simple Validation
func validatePath(path string) {
	keyword := strings.Split(path, "/")
	for _, value := range keyword {
		for _, character := range value {
			if !(character >= 65 && character <= 90 || character >= 97 && character <= 122) {
				panic(fmt.Sprintf("Invalid path %v", value))
			}
		}
	}
}

func (server *WebServer) Start() {
	address := server.ip.String() + ":" + strconv.Itoa(server.port)

	err := http.ListenAndServe(address, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		request = request.WithContext(server.ctx)
		server.router.ServeHTTP(writer, request)
	}))

	if err != nil {
		server.shutDown()
		return
	}

	// Do we need this part?
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	<-termChan

	server.shutDown()
}

func (server *WebServer) shutDown() {
	fmt.Println("The server is shutting Down... Will wait for 10 Second!")
	server.cancelFunc()
	time.Sleep(10 * time.Second)
}
