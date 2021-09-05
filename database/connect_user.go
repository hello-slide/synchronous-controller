package database

import (
	"fmt"
)

type DBConnectUsers struct {
	AbstractDBController
}

type ConnectUser struct {
	Id     string `db:"id"`
	UserId string `db:"user_id"`
}

// Create instance of user connect db.
//
// Arguments:
//	tableName {string} - table name.
//	config {Config} - db config.
//
// Returns:
//	{*DBConnectUsers} - user connect db instance.
func NewDBConnectUsers(tableName string, config *Config) (*DBConnectUsers, error) {
	db, err := NewDatabase(config)
	if err != nil {
		return nil, err
	}
	columns := "(id VARCHAR(256) NOT NULL, user_id VARCHAR(256) NOT NULL)"

	return &DBConnectUsers{
		AbstractDBController{
			DB:        db,
			TableName: tableName,
			Columns:   columns,
		},
	}, nil
}

// Add the participating users.
//
// Arguments:
//	data {ConnetUser} - user data.
func (c *DBConnectUsers) AddUser(data *ConnectUser) error {
	sql := fmt.Sprintf("INSERT INTO %s (id , user_id) VALUES ($1, $2)", c.TableName)

	_, err := c.DB.Execute(sql, data.Id, data.UserId)
	if err != nil {
		return err
	}
	return nil
}

// Gets the number of users participating in the target id.
//
// Arguments:
//	targetId {string} - Target id.
//
// Returns:
//	{int} - Number of participants.
func (c *DBConnectUsers) GetUserNumber(targetId string) (int, error) {
	sql := fmt.Sprintf("SELECT COUNT(id) FROM %s WHERE id = $1", c.TableName)

	return c.DB.Count(sql, targetId)
}

// Delete all target id information.
//
// Arguments:
//	targetId {string} - Target id to delete.
func (c *DBConnectUsers) Delete(targetId string) error {
	sql := fmt.Sprintf("DELETE FROM %s WHERE id = $1", c.TableName)

	_, err := c.DB.Execute(sql, targetId)
	if err != nil {
		return err
	}
	return nil
}
