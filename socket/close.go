package socket

import (
	"github.com/hello-slide/synchronous-controller/database"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

func Close(ws *websocket.Conn, db *database.DatabaseOp, id string, err error) {
	logrus.Infof("websocket connection failed err: %v", err)

	topic := database.NewDBTopic(TopicTableName, db)
	if err := topic.Delete(id); err != nil {
		logrus.Errorf("The topic could not be deleted. err: %v", err)
		return
	}

	answers := database.NewDBAnswers(AnswersTableName, db)
	if err := answers.Delete(id); err != nil {
		logrus.Errorf("The answers could not be deleted. err: %v", err)
		return
	}

	connectUsers := database.NewDBConnectUsers(ConnectUsersTablename, db)
	if err := connectUsers.Delete(id); err != nil {
		logrus.Errorf("The connect users could not be deleted. err: %v", err)
		return
	}

	db.Close()
	ws.Close()
}
