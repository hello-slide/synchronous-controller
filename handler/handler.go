package handler

import (
	"net/http"

	"golang.org/x/net/websocket"
)

func HostHandler(w http.ResponseWriter, r *http.Request) {
	s := websocket.Server{
		Handler: websocket.Handler(hostSocketHandler),
	}
	s.ServeHTTP(w, r)
}

func VisitorHandler(w http.ResponseWriter, r *http.Request) {
	s := websocket.Server{
		Handler: websocket.Handler(visitorSocketHandler),
	}
	s.ServeHTTP(w, r)
}
