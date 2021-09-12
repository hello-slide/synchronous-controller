package socket

import (
	"time"

	"github.com/hello-slide/synchronous-controller/database"
	"github.com/sirupsen/logrus"
)

func TopicGetter(db *database.DatabaseOp, topics *map[string]*string) {
	topic := database.NewDBTopic(TopicTableName, db)
	var updates = make(map[string]bool)

	for {
		for id, value := range *topics {
			exist, err := topic.Exist(id)
			if err != nil {
				logrus.Infof("topic exist err: %v", err)
				continue
			}
			if !exist {
				delete(*topics, id)
				continue
			}

			// update topic check
			newIsUpdate, err := topic.GetIsUpdate(id)
			if err != nil {
				logrus.Errorf("sendVisitor isUpdate error: %v", err)
				break
			}
			if isUpdate, ok := updates[id]; ok {
				if newIsUpdate != isUpdate {
					topicData, err := topic.GetTopic(id)
					if err != nil {
						logrus.Errorf("get topic err: %v", err)
						continue
					}
					*value = topicData

					updates[id] = newIsUpdate
				}
			}else {
				topicData, err := topic.GetTopic(id)
					if err != nil {
						logrus.Errorf("get topic err: %v", err)
						continue
					}
				*value = topicData

				updates[id] = newIsUpdate
			}
		}

		// delete extra buffer flags.
		for id := range updates {
			var exist = false
			for existId := range *topics {
				if existId == id {
					exist = true
				}
			}

			if !exist {
				delete(updates, id)
			}
		}

		time.Sleep(1 * time.Second)
	}
}
