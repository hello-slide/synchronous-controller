package socket

import (
	"github.com/hello-slide/synchronous-controller/database"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

type CloseSocket struct {
	ws *websocket.Conn
	db *database.DatabaseOp
	id string
}

// Create close socket instance.
//
// Arguments:
//	ws {*websocket.Conn} - websocket
//	db {*database.DatabaseOp} - database operator
//	id {string} - id.
func NewCloseSocket(ws *websocket.Conn, db *database.DatabaseOp, id string) *CloseSocket {
	return &CloseSocket{
		ws: ws,
		db: db,
		id: id,
	}
}

func (c *CloseSocket) HostNoErr() {
	if err := c.Host(); err != nil {
		logrus.Infof("ERR socket close: %v", err)
	}
}

func (c *CloseSocket) VisitorNoErr(userId string, queue *map[string]map[string]*websocket.Conn) {
	if err := c.Visitor(userId, queue); err != nil {
		logrus.Infof("ERR socket close: %v", err)
	}
}

// Close host.
//
//	- delete topic for id.
//	- delete answers for id.
//	- delete connect users for id.
func (c *CloseSocket) Host() error {
	defer c.ws.Close()

	logrus.Infof("close socket id: %v", c.id)

	topic := database.NewDBTopic(TopicTableName, c.db)
	if err := topic.Delete(c.id); err != nil {
		return err
	}

	answers := database.NewDBAnswers(AnswersTableName, c.db)
	if err := answers.Delete(c.id); err != nil {
		return err
	}

	connectUsers := database.NewDBConnectUsers(ConnectUsersTablename, c.db)
	if err := connectUsers.Delete(c.id); err != nil {
		return err
	}
	return nil
}

// Close visitor.
//
//	- delete  connect users for user id.
func (c *CloseSocket) Visitor(userId string, queue *map[string]map[string]*websocket.Conn) error {
	delete((*queue)[c.id], userId)

	defer c.ws.Close()

	logrus.Infof("close socket visitor id: %v, userid: %v", c.id, userId)

	connectUser := database.NewDBConnectUsers(ConnectUsersTablename, c.db)
	if err := connectUser.DeleteUser(userId); err != nil {
		return err
	}
	return nil
}
