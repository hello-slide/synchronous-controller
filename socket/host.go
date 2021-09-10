package socket

import (
	"strconv"
	"time"

	"github.com/hello-slide/synchronous-controller/database"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

type SendAnswers struct {
	SendType string            `json:"type"`
	Answers  []database.Answer `json:"answers"`
}

func SendHost(ws *websocket.Conn, quit chan bool, db *database.DatabaseOp, id string) {
	connectUser := database.NewDBConnectUsers(ConnectUsersTablename, db)
	answers := database.NewDBAnswers(AnswersTableName, db)

	var usersBuffer int = 0
	var answersBuffer []database.Answer = []database.Answer{}

	for {
		select {
		case <-quit:
			return
		default:
			// participant
			nums, err := connectUser.GetUserNumber(id)
			if err != nil {
				logrus.Errorf("sendHost getusernumber error: %v", err)
				return
			}
			if nums != usersBuffer {
				sendData := map[string]string{
					"type":     "2",
					"visitors": strconv.Itoa(nums),
				}
				if err := websocket.JSON.Send(ws, sendData); err != nil {
					logrus.Errorf("sendHost send visitors error: %v", err)
					return
				}

				usersBuffer = nums
			}

			time.Sleep(1 * time.Second)

			// answers
			_answers, err := answers.GetAnswers(id)
			if err != nil {
				logrus.Errorf("sendHost getanswers error: %v", err)
				return
			}

			if len(_answers) != len(answersBuffer) {
				sendData := &SendAnswers{
					SendType: "3",
					Answers:  _answers,
				}
				if err := websocket.JSON.Send(ws, sendData); err != nil {
					logrus.Errorf("sendHost send answers error: %v", err)
					return
				}

				answersBuffer = _answers
			}

			time.Sleep(1 * time.Second)
		}
	}
}

func ReceiveHost(ws *websocket.Conn, quit chan bool, db *database.DatabaseOp, id string) {
	topic := database.NewDBTopic(TopicTableName, db)
	isUpdate, err := topic.GetIsUpdate(id)
	if err != nil {
		logrus.Errorf("receivedHost isUpdate error: %v", err)
		return
	}

	for {
		select {
		case <-quit:
			return
		default:
			var receivedData map[string]string
			websocket.JSON.Receive(ws, receivedData)

			statusType, ok1 := receivedData["type"]
			newTopic, ok2 := receivedData["topic"]

			if ok1 && ok2 && statusType == "4" {
				topicData := &database.Topic{
					Id:       id,
					IsUpdate: !isUpdate,
					Topic:    newTopic,
				}
				if err := topic.UpdateTopic(topicData); err != nil {
					logrus.Errorf("receivedHost updateTopic error: %v", err)
					return
				}
				isUpdate = !isUpdate
			}

			time.Sleep(1 * time.Second)
		}
	}
}
