package socket

import (
	"github.com/hello-slide/synchronous-controller/database"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

// Close operation.
//
// close the all database connection and websocket.
//
// Arguments:
//	ws {*websocket.Conn} - websocket op.
//	db {*database.DatabaseOp} - database operator.
//	id {string} - session id or user id.
func Close(ws *websocket.Conn, db *database.DatabaseOp, status Status, id string) {
	defer ws.Close()

	if status == Host {
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
	} else if status == Visitor {
		connectUser := database.NewDBConnectUsers(ConnectUsersTablename, db)
		if err := connectUser.DeleteUser(id); err != nil {
			logrus.Errorf("The connect users could not be deleted. err: %v", err)
			return
		}
	}

	db.Close()
}
