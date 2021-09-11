package socket

import (
	"io"
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

// Send a socket to the host.
//
// Arguments:
//	ws {*websocket.Conn} - websocket operator.
//	db {*database.DatabaseOp} - database op.
//	id {string} - id
//	quit {chan bool} - quit signal.
func SendHost(ws *websocket.Conn, db *database.DatabaseOp, id string, quit chan bool) {
	connectUser := database.NewDBConnectUsers(ConnectUsersTablename, db)
	answers := database.NewDBAnswers(AnswersTableName, db)

	var usersBuffer int = 0
	var answersBuffer []string = []string{}

	for {
		select {
		case <- quit:
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
				sendAns := choice(&_answers, &answersBuffer)

				sendData := &SendAnswers{
					SendType: "3",
					Answers:  sendAns,
				}
				if err := websocket.JSON.Send(ws, sendData); err != nil {
					logrus.Errorf("sendHost send answers error: %v", err)
					return
				}
			}
			time.Sleep(1 * time.Second)
		}
	}
}

// select answers.
//
// Arguments:
//	answers {[]database.Answer} - answers.
//	buffer {[]string} - buffer list
//
// Returns:
//	{[]database.Answer} - send ansers.
func choice(answers *[]database.Answer, buffer *[]string) ([]database.Answer) {
	sendAns := []database.Answer{}

	// reset answers buffer
	if len(*answers) == 0 {
		*buffer = []string{}
		return *answers
	}

	for _, ans := range *answers {
		isExist := false
		for _, buf := range *buffer {
			if buf == ans.UserId {
				isExist = true
				break
			}
		}

		if !isExist {
			*buffer = append(*buffer, ans.UserId)
			sendAns = append(sendAns, ans)
		}
	}

	return sendAns
}

// Received a socket to the host.
//
// Arguments:
//	ws {*websocket.Conn} - websocket operator.
//	db {*database.DatabaseOp} - database op.
//	id {string} - id
//	quit {chan bool} - quit signal.
func ReceiveHost(ws *websocket.Conn, db *database.DatabaseOp, id string, quit chan bool) {
	topic := database.NewDBTopic(TopicTableName, db)
	answers := database.NewDBAnswers(AnswersTableName, db)

	isUpdate, err := topic.GetIsUpdate(id)
	if err != nil {
		logrus.Errorf("receivedHost isUpdate error: %v", err)
		return
	}

	for {
		var receivedData map[string]string
		if err := websocket.JSON.Receive(ws, &receivedData); err != nil {
			if err == io.EOF {
				quit <- true
			}else{
				logrus.Errorf("websocket recrived error: %v", err)
			}
			return
		}

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

			// Sleep for a period of time as they may not be deleted if answered after updating the topic.
			time.Sleep(2 * time.Second)

			if err := answers.Delete(id); err != nil {
				logrus.Errorf("receivedHost deleteAnswers error: %v", err)
				return
			}
			isUpdate = !isUpdate
		}
	}

}
