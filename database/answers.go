package database

import "fmt"

type Answer struct {
	Id     string `db:"id" json:"id"`
	UserId string `db:"user_id" json:"user_id"`
	Name   string `db:"name" json:"name"`
	Answer string `db:"answer" json:"answer"`
}

type DBAnswers struct {
	AbstractDBController
}

// Create instance of answers db.
//
// Arguments:
//	tableName {string} - table name.
//	db {*DatabaseOp} - database instance.
//
// Returns:
//	{*DBConnectUsers} - user connect db instance.
func NewDBAnswers(tableName string, db *DatabaseOp) *DBAnswers {
	columns := "(id VARCHAR(256) NOT NULL, user_id VARCHAR(256) NOT NULL, name VARCHAR(256), answer VARCHAR(1024), PRIMARY KEY (user_id))"

	return &DBAnswers{
		AbstractDBController{
			DB:        db,
			TableName: tableName,
			Columns:   columns,
		},
	}
}

// Insert the new answer.
//
// Arguments:
// data {*Answer} - Answer data.
func (c *DBAnswers) AddAnswer(data *Answer) error {
	sql := fmt.Sprintf("INSERT INTO %s (id , user_id, name, answer) VALUES ($1, $2, $3, $4)", c.TableName)

	_, err := c.DB.Execute(sql, data.Id, data.UserId, data.Name, data.Answer)
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
		if err := result.Scan(&data.Id, &data.UserId, &data.Name, &data.Answer); err != nil {
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
