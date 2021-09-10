package handler

import (
	"github.com/google/uuid"
	"github.com/hello-slide/synchronous-controller/socket"
	"github.com/hello-slide/synchronous-controller/util"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

// websocket handler of host.
func hostSocketHandler(ws *websocket.Conn) {
	logrus.Debugf("Create host socket.")

	id, err := socket.Init(ws, socket.Host, db, "")
	if err != nil {
		logrus.Infof("socket error: %v", err)
		ws.Close()
		return
	}
	defer socket.Close(ws, db, socket.Host, id)

	socket.SendHost(ws, db, id)
	go socket.ReceiveHost(ws, db, id)
}

// websocket handler of visitor.
func visitorSocketHandler(ws *websocket.Conn) {
	logrus.Debugf("Create visitor socket.")

	uuidObj, err := uuid.NewUUID()
	if err != nil {
		logrus.Infof("uuid error: %v", err)
		ws.Close()
		return
	}
	userId := util.NewDateSeed().AddSeed(uuidObj.String()).CreateSpecifyLength(5)

	id, err := socket.Init(ws, socket.Visitor, db, userId)
	if err != nil {
		logrus.Infof("socket error: %v", err)
		ws.Close()
		return
	}

	defer socket.Close(ws, db, socket.Visitor, userId)

	socket.SendVisitor(ws, db, id)
	go socket.ReceiveVisitor(ws, db, id, userId)
}
