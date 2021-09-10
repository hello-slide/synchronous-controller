package handler

import (
	"net/http"

	"golang.org/x/net/websocket"
)

func Roothandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

// handler of host.
func HostHandler(w http.ResponseWriter, r *http.Request) {
	s := websocket.Server{
		Handler: websocket.Handler(hostSocketHandler),
	}
	s.ServeHTTP(w, r)
}

// handler of visitor.
func VisitorHandler(w http.ResponseWriter, r *http.Request) {
	s := websocket.Server{
		Handler: websocket.Handler(visitorSocketHandler),
	}
	s.ServeHTTP(w, r)
}
