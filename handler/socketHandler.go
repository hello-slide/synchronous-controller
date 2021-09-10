package handler

import (
	"github.com/hello-slide/synchronous-controller/database"
	"github.com/hello-slide/synchronous-controller/socket"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

// websocket handler of host.
func hostSocketHandler(ws *websocket.Conn) {
	db, err := database.NewDatabase(dbConfig)
	if err != nil {
		logrus.Infof("connect db error: %v", err)
		ws.Close()
		return
	}

	id, err := socket.Init(ws, socket.Host, db)
	if err != nil {
		logrus.Infof("socket error: %v", err)
		ws.Close()
		return
	}
	defer socket.Close(ws, db, id)

	socket.SendHost(ws, db, id)
	go socket.ReceiveHost(ws, db, id)
}

// websocket handler of visitor.
func visitorSocketHandler(ws *websocket.Conn) {
	defer ws.Close()
}
