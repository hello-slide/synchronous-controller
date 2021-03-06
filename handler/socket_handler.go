package handler

import (
	"io"

	"github.com/hello-slide/synchronous-controller/socket"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

var topics = make(map[string]*string)

// websocket handler of host.
func hostSocketHandler(ws *websocket.Conn) {
	id, err := socket.NewInitSocket(ws, db).Host()
	if err != nil {
		if err != io.EOF {
			logrus.Errorf("host init socket err: %v", err)
		}
		ws.Close()
		return
	}

	defer socket.NewCloseSocket(ws, db, id).HostNoErr()

	quit := make(chan bool)

	go socket.SendHost(ws, db, id, quit)
	socket.ReceiveHost(ws, db, id, quit)
}

// websocket handler of visitor.
func visitorSocketHandler(ws *websocket.Conn) {
	id, userId, err := socket.NewInitSocket(ws, db).Visitor()
	if err != nil {
		if err != io.EOF {
			logrus.Errorf("visitor init socket err: %v", err)
		}
		ws.Close()
		return
	}

	defer socket.NewCloseSocket(ws, db, id).VisitorNoErr(userId)

	quit := make(chan bool)

	go socket.SendVisitor(ws, db, id, userId, quit, &topics)
	socket.ReceiveVisitor(ws, db, id, userId, quit)
}

func BackgroundTopicGetter() {
	socket.TopicGetter(db, &topics)
}
