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
				endWebsocket(&element)
				return
			}
			if !exist {
				endWebsocket(&element)
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


func endWebsocket(sockets *map[string]*websocket.Conn) {
	for _, ws := range *sockets {
		ws.Close()
	}
}

func sendTopics(topics *map[string]*websocket.Conn, topic *database.DBTopic, id string) {
	for _, ws := range *topics {
		if err := sendTopic(ws, topic, id); err != nil {
			logrus.Errorf("send topic error: %v", err)
			return
		}
	}
}

// send topics.
//
// Arguments:
//	ws {*websocket.Conn} - websocket operator.
//	db {*database.DatabaseOp} - database op.
//	id {string} - id
func sendTopic(ws *websocket.Conn, topic *database.DBTopic, id string) error {
	topicData, err := topic.GetTopic(id)
	if err != nil {
		return err
	}
	sendData := map[string]string{
		"type":  "5",
		"topic": topicData,
	}
	if err := websocket.JSON.Send(ws, sendData); err != nil {
		return err
	}

	return nil
}
