package database

import "fmt"

type Answer struct {
	Id     string `db:"id"`
	UserId string `db:"user_id"`
	Answer string `db:"answer"`
}

type DBAnswers struct {
	AbstractDBController
}

// Create instance of answers db.
//
// Arguments:
//	tableName {string} - table name.
//	config {Config} - db config.
//
// Returns:
//	{*DBConnectUsers} - user connect db instance.
func NewDBAnswers(tableName string, config *Config) (*DBAnswers, error) {
	db, err := NewDatabase(config)
	if err != nil {
		return nil, err
	}
	columns := "(id VARCHAR(256) NOT NULL, user_id VARCHAR(256) NOT NULL, answer VARCHAR(1024))"

	return &DBAnswers{
		AbstractDBController{
			DB:        db,
			TableName: tableName,
			Columns:   columns,
		},
	}, nil
}

// Insert the new answer.
//
// Arguments:
// data {*Answer} - Answer data.
func (c *DBAnswers) AddAnswer(data *Answer) error {
	sql := fmt.Sprintf("INSERT INTO %s (id , user_id, answer) VALUES ($1, $2, $3)", c.TableName)

	_, err := c.DB.Execute(sql, data.Id, data.UserId, data.Answer)
	if err != nil {
		return err
	}
	return nil
}

// Get the all answers.
//
// Arguments:
//	targetId {string} - target id.
//
// Returns:
//	{[]Answer} - all answers.
func (c *DBAnswers) GetAnswers(targetId string) ([]Answer, error) {
	sql := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", c.TableName)
	result, err := c.DB.Query(sql, targetId)
	if err != nil {
		return nil, err
	}

	var answers []Answer = []Answer{}
	var data Answer

	for result.Next() {
		if err := result.Scan(&data.Id, &data.UserId, &data.Answer); err != nil {
			return nil, err
		}
		answers = append(answers, data)
	}
	if err := result.Err(); err != nil {
		return nil, err
	}
	result.Close()

	return answers, nil
}

// Delete all target id information.
//
// Arguments:
//	targetId {string} - Target id to delete.
func (c *DBAnswers) Delete(targetId string) error {
	sql := fmt.Sprintf("DELETE FROM %s WHERE id = $1", c.TableName)

	_, err := c.DB.Execute(sql, targetId)
	if err != nil {
		return err
	}
	return nil
}
