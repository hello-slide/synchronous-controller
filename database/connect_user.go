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
func NewDBConnectUsers(tableName string, config Config) (*DBConnectUsers, error) {
	db, err := NewDatabase(config)
	if err != nil {
		return nil, err
	}

	return &DBConnectUsers{
		AbstractDBController{
			DB:        db,
			TableName: tableName,
		},
	}, nil
}

// Add the participating users.
//
// Arguments:
//	data {ConnetUser} - user data.
func (c *DBConnectUsers) AddUser(data ConnectUser) error {
	sql := "INSERT INTO ? (id , user_id) NOT NULL) VALUES (?, ?)"

	_, err := c.DB.Execute(sql, c.TableName, data.Id, data.UserId)
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
	sql := "SELECT COUNT(id) FROM ? WHERE id = '?'"

	var count interface{}

	if err := c.DB.QueryOneRecord(&count, sql, c.TableName, targetId); err != nil {
		return 0, err
	}

	if result, ok := count.(int); ok {
		return result, nil
	}
	return 0, fmt.Errorf("the result could not be parsed to empty or int type result: %v", count)
}

// Delete all target id information.
//
// Arguments:
//	targetId {string} - Target id to delete.
func (c *DBConnectUsers) Delete(targetId string) error {
	sql := "DELETE FROM ? WHERE id = '?'"

	_, err := c.DB.Execute(sql, c.TableName, targetId)
	if err != nil {
		return err
	}
	return nil
}
