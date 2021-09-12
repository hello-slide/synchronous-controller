package socket

import (
	"io"

	"github.com/hello-slide/synchronous-controller/database"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

// Send a socket to the visitor.
//
// Arguments:
//	ws {*websocket.Conn} - websocket operator.
//	db {*database.DatabaseOp} - database op.
//	id {string} - id
//	userId {string} - user id
//	quit {chan bool} - quit signal.
//	queue {*map[string]map[string]*websocket.Conn} - send visitor queue.
func SendVisitor(ws *websocket.Conn, db *database.DatabaseOp, id string, userId string, quit chan bool, queue *map[string]*map[string]*websocket.Conn) {
	select {
	case <- quit:
		return
	default:
		topic := database.NewDBTopic(TopicTableName, db)

		exist, err := topic.Exist(id)
		if err != nil {
			ws.Close()
			return
		}
		if !exist {
			ws.Close()
			return
		}

		if err := sendTopic(ws, topic, id); err != nil {
			logrus.Errorf("error: %v", err)
			return
		}

		(*(*queue)[id])[userId] = ws
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


// Received a socket to the visitor.
//
// Arguments:
//	ws {*websocket.Conn} - websocket operator.
//	db {*database.DatabaseOp} - database op.
//	id {string} - id
//	quit {chan bool} - quit signal.
func ReceiveVisitor(ws *websocket.Conn, db *database.DatabaseOp, id string, userId string, quit chan bool) {
	answers := database.NewDBAnswers(AnswersTableName, db)
	for {
		var receivedData map[string]string
		if err := websocket.JSON.Receive(ws, &receivedData); err != nil {
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
