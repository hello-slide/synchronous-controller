package handler

import (
	"github.com/hello-slide/synchronous-controller/database"
	"github.com/hello-slide/synchronous-controller/socket"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

func hostSocketHandler(ws *websocket.Conn) {
	db, err := database.NewDatabase(dbConfig)
	if err != nil {
		logrus.Infof("connect db error: %v", err)
		return
	}

	id, err := socket.Init(ws, socket.Host, db)
	if err != nil {
		logrus.Infof("socket error: %v", err)
		return
	}
	defer socket.Close(ws, db, id, nil)

	socket.SendHost(ws, id)

}

func visitorSocketHandler(ws *websocket.Conn) {
	defer ws.Close()
}
