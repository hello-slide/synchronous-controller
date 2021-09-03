package database

import (
	"fmt"
)

type DBConnectUsers struct {
	DB        *DatabaseOp
	TableName string
}

type ConnectUser struct {
	Id     string `db:"id"`
	UserId string `db:"user_id"`
}

func NewDBConnectUsers(tableName string, config Config) (*DBConnectUsers, error) {
	db, err := NewDatabase(config)
	if err != nil {
		return nil, err
	}

	return &DBConnectUsers{
		DB:        db,
		TableName: tableName,
	}, nil
}

func (c *DBConnectUsers) CreateTable() error {
	columns := "(id VARCHAR(256) NOT NULL, user_id VARCHAR(256) NOT NULL)"

	return c.DB.CreateTable(c.TableName, columns)
}

func (c *DBConnectUsers) AddUser(data ConnectUser) error {
	sql := "INSERT INTO ? (id , user_id) NOT NULL) VALUES (?, ?)"

	_, err := c.DB.Execute(sql, c.TableName, data.Id, data.UserId)
	if err != nil {
		return err
	}
	return nil
}

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

func (c *DBConnectUsers) Delete(targetId string) error {
	sql := "DELETE FROM ? WHERE id = '?'"

	_, err := c.DB.Execute(sql, c.TableName, targetId)
	if err != nil {
		return err
	}
	return nil
}
