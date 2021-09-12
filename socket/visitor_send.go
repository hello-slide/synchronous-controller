package socket

import (
	"time"

	"github.com/hello-slide/synchronous-controller/database"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

func VisitorSend(db *database.DatabaseOp, queue *map[string]map[string]*websocket.Conn) {
	topic := database.NewDBTopic(TopicTableName, db)

	var updates = make(map[string]bool)

	for {
		for id, element := range *queue {
			if element == nil {
				continue
			}

			exist, err := topic.Exist(id)
			if err != nil {
				endWebsocket(&element, id)
				return
			}
			if !exist {
				endWebsocket(&element, id)
				return
			}

			newIsUpdate, err := topic.GetIsUpdate(id)
			if err != nil {
				logrus.Errorf("sendVisitor isUpdate error: %v", err)
				return
			}

			if isUpdate, ok := updates[id]; ok {
				if newIsUpdate != isUpdate {
					sendTopics(&element, topic, id)

					updates[id] = newIsUpdate
				}
			}else {
				updates[id] = newIsUpdate
			}
		}

		time.Sleep(1 * time.Second)
	}
}

// Close websocket.
func endWebsocket(sockets *map[string]*websocket.Conn, id string) {
	logrus.Infof("close visitors id: %v", id)

	for _, ws := range *sockets {
		ws.Close()
	}
}

// send topic to ids.
func sendTopics(conns *map[string]*websocket.Conn, topic *database.DBTopic, id string) {
	topicData, err := topic.GetTopic(id)
	if err != nil {
		logrus.Errorf("get topic err: %v", err)
		return
	}
	sendData := map[string]string{
		"type":  "5",
		"topic": topicData,
	}

	logrus.Infof("update topic: %v", topicData)

	for _, ws := range *conns {
		if err := websocket.JSON.Send(ws, sendData); err != nil {
			logrus.Errorf("websocket send err: %v", err)
		}
	}
}
