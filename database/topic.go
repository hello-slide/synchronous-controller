package database

import "fmt"

type DBTopic struct {
	AbstractDBController
}

type Topic struct {
	Id       string `db:"id"`
	Topic    string `db:"topic"`
	IsUpdate bool   `db:"is_update"`
}

// Create instance of topic db.
//
// Arguments:
//	tableName {string} - table name.
//	config {Config} - db config.
//
// Returns:
//	{*DBConnectUsers} - topic db instance.
func NewDBTopic(tableName string, config *Config) (*DBTopic, error) {
	db, err := NewDatabase(config)
	if err != nil {
		return nil, err
	}
	columns := "(id VARCHAR(256) NOT NULL,is_update boolean, topic VARCHAR(1024))"

	return &DBTopic{
		AbstractDBController{
			DB:        db,
			TableName: tableName,
			Columns:   columns,
		},
	}, nil
}

// Create topics.
//
// Arguments:
//	topic {*Topic} - topic data.
func (c *DBTopic) CreateTopic(topic *Topic) error {
	sql := fmt.Sprintf("INSERT INTO %s (id , is_update, topic) VALUES ($1, false, $2)", c.TableName)

	_, err := c.DB.Execute(sql, topic.Id, topic.Topic)
	if err != nil {
		return err
	}
	return nil
}

// Update topics.
//
// Arguments:
//	topic {*Topic} - new topic
func (c *DBTopic) UpdateTopic(topic *Topic) error {
	sql := fmt.Sprintf("UPDATE %s SET topic = $1, is_update = $2 WHERE id = $3", c.TableName)

	_, err := c.DB.Execute(sql, topic.Topic, topic.IsUpdate, topic.Id)
	if err != nil {
		return err
	}
	return nil
}

// Get the topic.
//
// Arguments:
//	id {string} - topic id.
//
// Returns:
//	{string} - topic text.
func (c *DBTopic) GetTopic(id string) (string, error) {
	sql := fmt.Sprintf("SELECT topic FROM %s WHERE id = $1", c.TableName)

	result, err := c.DB.Query(sql, id)
	if err != nil {
		return "", err
	}

	result.Next()
	var topic string

	if err := result.Scan(&topic); err != nil {
		return "", err
	}
	if err := result.Err(); err != nil {
		return "", err
	}
	result.Close()

	return topic, nil

}

// Get update flag.
//
// Arguments:
//	id {string} - topic id.
//
// Returns:
//	{bool} - update flag.
func (c *DBTopic) GetIsUpdate(id string) (bool, error) {
	sql := fmt.Sprintf("SELECT is_update FROM %s WHERE id = $1", c.TableName)

	result, err := c.DB.Query(sql, id)
	if err != nil {
		return false, err
	}

	result.Next()
	var isUpdate bool

	if err := result.Scan(&isUpdate); err != nil {
		return false, err
	}
	if err := result.Err(); err != nil {
		return false, err
	}

	return isUpdate, nil
}
