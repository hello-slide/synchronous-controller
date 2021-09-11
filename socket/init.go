package socket

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/hello-slide/synchronous-controller/database"
	"github.com/hello-slide/synchronous-controller/util"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

type Status int

const (
	Host    Status = iota
	Visitor        = iota
)

type InitializeSocket struct {
	ws *websocket.Conn
	db *database.DatabaseOp
}

// Create initialize socket instance.
//
// Arguments:
//	ws {*websocket.Conn} - websocket
//	db {*database.DatabaseOp} - database operator
func NewInitSocket(ws *websocket.Conn, db *database.DatabaseOp) *InitializeSocket {
	return &InitializeSocket{
		ws: ws,
		db: db,
	}
}

// initialize host.
//
// Returns:
//	{string} - id
func (c *InitializeSocket) Host() (string, error) {
	var responseMessage map[string]string
	if err := websocket.JSON.Receive(c.ws, &responseMessage); err != nil {
		return "", err
	}

	logrus.Debugf("Init response data: %v", responseMessage)

	responseType, ok := responseMessage["type"]
	if !ok {
		return "", errors.New("you need to specify the type")
	}
	if responseType != "0" {
		return "", errors.New("the initial type must be 0")
	}

	connectUser := database.NewDBConnectUsers(ConnectUsersTablename, c.db)

	// host
	id, err := c.createId("")
	if err != nil {
		return "", err
	}

	// Databas operation
	answers := database.NewDBAnswers(AnswersTableName, c.db)
	if err := answers.CreateTable(); err != nil {
		return "", err
	}
	if err := connectUser.CreateTable(); err != nil {
		return "", err
	}

	topicOp := database.NewDBTopic(TopicTableName, c.db)
	if err := topicOp.CreateTable(); err != nil {
		return "", err
	}
	topic := &database.Topic{
		Id:       id,
		Topic:    "",
		IsUpdate: false,
	}
	if err := topicOp.CreateTopic(topic); err != nil {
		return "", err
	}

	initializeSendMessage := map[string]string{
		"type":    "0",
		"version": "1.0",
		"id":      id,
	}

	if err := websocket.JSON.Send(c.ws, initializeSendMessage); err != nil {
		return "", err
	}

	return id, nil
}

// initialize visitor
//
// Returns:
//	{string} - id.
//	{string} - user id.
func (c *InitializeSocket) Visitor() (string, string, error) {
	var responseMessage map[string]string
	if err := websocket.JSON.Receive(c.ws, &responseMessage); err != nil {
		return "", "", err
	}

	logrus.Debugf("Init response data: %v", responseMessage)

	responseType, ok := responseMessage["type"]
	if !ok {
		return "", "", errors.New("you need to specify the type")
	}
	if responseType != "1" {
		return "", "", errors.New("the initial type must be 1")
	}

	connectUser := database.NewDBConnectUsers(ConnectUsersTablename, c.db)
	topic := database.NewDBTopic(TopicTableName, c.db)

	id, ok := responseMessage["id"]
	if !ok {
		return "", "", errors.New("id is not found")
	}
	isExist, err := topic.Exist(id)
	if err != nil {
		return "", "", err
	}
	if !isExist {
		return "", "", errors.New("id is not exists")
	}

	userId, err := c.createId(id)
	if err != nil {
		return "", "", err
	}

	if err := connectUser.AddUser(&database.ConnectUser{
		Id:     id,
		UserId: userId,
	}); err != nil {
		return "", "", err
	}

	// visitor
	initializeSendMessage := map[string]string{
		"type":    "1",
		"version": "1.0",
	}

	if err := websocket.JSON.Send(c.ws, initializeSendMessage); err != nil {
		return "", "", err
	}

	return id, userId, nil
}

// Create id.
//
// Arguments:
//	prefix {string} - id prefix.
//
// Returns:
//	{string} - id.
func (c *InitializeSocket) createId(prefix string) (string, error) {
	uuidObj, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	id := util.NewDateSeed().AddSeed(uuidObj.String()).CreateSpecifyLength(5)

	if len(prefix) == 0 {
		return id, nil
	}

	return strings.Join([]string{prefix, id}, ""), nil
}
