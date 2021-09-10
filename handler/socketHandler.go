package handler

import (
	"time"

	"github.com/hello-slide/synchronous-controller/database"
	"github.com/hello-slide/synchronous-controller/socket"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

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
	defer socket.Close(ws, db, id, nil)

	quit := make(chan bool)

	socket.SendHost(ws, quit, db, id)
	go socket.ReceiveHost(ws, quit, db, id)

	time.Sleep(2 * time.Second)
	quit <- true
}

func visitorSocketHandler(ws *websocket.Conn) {
	defer ws.Close()
}
