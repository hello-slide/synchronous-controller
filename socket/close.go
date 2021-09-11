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

func (c *CloseSocket) VisitorNoErr() {
	if err := c.Visitor(); err != nil {
		logrus.Infof("ERR socket close: %v", err)
	}
}

func (c *CloseSocket) Host() error {
	defer c.ws.Close()

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

func (c *CloseSocket) Visitor() error {
	defer c.ws.Close()

	connectUser := database.NewDBConnectUsers(ConnectUsersTablename, c.db)
	if err := connectUser.DeleteUser(c.id); err != nil {
		return err
	}
	return nil
}
