package socket

import (
	"io"
	"time"

	"github.com/hello-slide/synchronous-controller/database"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

func SendVisitor(ws *websocket.Conn, db *database.DatabaseOp, id string, quit chan bool) {
	topic := database.NewDBTopic(TopicTableName, db)

	if err := sendTopic(ws, topic, id); err != nil {
		logrus.Errorf("error: %v", err)
		return
	}
	isUpdate, err := topic.GetIsUpdate(id)
	if err != nil {
		logrus.Errorf("sendVisitor isUpdate error: %v", err)
		return
	}

	for {
		select {
		case <- quit:
			return
		default:
			newIsUpdate, err := topic.GetIsUpdate(id)
			if err != nil {
				logrus.Errorf("sendVisitor isUpdate error: %v", err)
				return
			}

			if newIsUpdate != isUpdate {
				if err := sendTopic(ws, topic, id); err != nil {
					logrus.Errorf("error: %v", err)
					return
				}

				isUpdate = newIsUpdate
			}
			time.Sleep(1 * time.Second)
		}
	}
}

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

func ReceiveVisitor(ws *websocket.Conn, db *database.DatabaseOp, id string, userId string, quit chan bool) {
	answers := database.NewDBAnswers(AnswersTableName, db)
	for {
		var receivedData map[string]string
		if err := websocket.JSON.Receive(ws, receivedData); err != nil {
			if err == io.EOF {
				quit <- true
				logrus.Infof("close socket id: %v", id)
			}else{
				logrus.Errorf("websocket recrived error: %v", err)
			}
			return
		}

		statusType, ok := receivedData["type"]

		if ok && statusType == "6" {
			data := &database.Answer{
				Id:     id,
				UserId: userId,
				Name:   receivedData["name"],
				Answer: receivedData["answer"],
			}

			if err := answers.AddAnswer(data); err != nil {
				logrus.Errorf("error: %v", err)
				return
			}
		}
	}
}
