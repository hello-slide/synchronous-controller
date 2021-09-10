package socket

import (
	"errors"

	"github.com/google/uuid"
	"github.com/hello-slide/synchronous-controller/database"
	"github.com/hello-slide/synchronous-controller/util"
	"golang.org/x/net/websocket"
)

type Status int

const (
	Host    Status = iota
	Visitor        = iota
)

// Initialize websocket connection.
//
// Assign an id on the Host.
//
// Arguments:
//	we {*websocket.Conn} - websocket conn.
//	status {Status} - Type to initialize. host or visitor.
//	db {*database.DatabaseOp} - database op instance.
//	userId {string} - user id. use only visitor.
//
// Returns:
//	{string} - unique id.
func Init(ws *websocket.Conn, status Status, db *database.DatabaseOp, userId string) (string, error) {
	var responseMessage map[string]string
	if err := websocket.JSON.Receive(ws, &responseMessage); err != nil {
		return "", err
	}

	responseType, ok := responseMessage["type"]
	if !ok {
		return "", errors.New("you need to specify the type")
	}

	connectUser := database.NewDBConnectUsers(ConnectUsersTablename, db)

	if responseType == "0" && status == Host {
		// host
		uuidObj, err := uuid.NewUUID()
		if err != nil {
			return "", err
		}
		token := util.NewDateSeed().AddSeed(uuidObj.String())
		id := token.CreateSpecifyLength(5)

		// Databas operation
		answers := database.NewDBAnswers(AnswersTableName, db)
		if err := answers.CreateTable(); err != nil {
			return "", nil
		}
		if err := connectUser.CreateTable(); err != nil {
			return "", nil
		}

		topicOp := database.NewDBTopic(TopicTableName, db)
		if err := topicOp.CreateTable(); err != nil {
			return "", err
		}
		topic := &database.Topic{
			Id:       id,
			Topic:    "",
			IsUpdate: false,
		}
		if err := topicOp.CreateTopic(topic); err != nil {
			return "", nil
		}

		initializeSendMessage := map[string]string{
			"type":    "0",
			"version": "1.0",
			"id":      id,
		}

		if err := websocket.JSON.Send(ws, initializeSendMessage); err != nil {
			return "", err
		}

		return id, nil
	}
	if responseType == "1" && status == Visitor {
		id, ok := responseMessage["id"]
		if !ok {
			return "", errors.New("id is not found")
		}

		if err := connectUser.AddUser(&database.ConnectUser{
			Id:     id,
			UserId: userId,
		}); err != nil {
			return "", err
		}

		// visitor
		initializeSendMessage := map[string]string{
			"type":    "1",
			"version": "1.0",
		}

		if err := websocket.JSON.Send(ws, initializeSendMessage); err != nil {
			return id, err
		}

		return "", nil
	}
	return "", errors.New("the initial type must be 0 or 1")
}
