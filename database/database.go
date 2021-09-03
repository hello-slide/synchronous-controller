package database

import (
	"database/sql"
	"fmt"
)

type DatabaseOp struct {
	Db *sql.DB
}

func NewDatabase(config Config) (*DatabaseOp, error) {
	db, err := sql.Open(config.driverName, config.dataSourceName)
	if err != nil {
		return nil, err
	}

	return &DatabaseOp{
		Db: db,
	}, nil
}

// SQL Query.
//
// Arguments:
//	sql {string} - text for sql.
//	args {...interface{}} - args
//
// Returns:
//	{*sql.Rows} - results.
func (c *DatabaseOp) Query(sql string, args ...interface{}) (*sql.Rows, error) {
	return c.Db.Query(sql, args...)
}

// SQL Query Row.
//
// Arguments:
//	dest {*interface{}} - dest value.
//	sql {string} - text for sql.
//	args {...interface{}} - args
func (c *DatabaseOp) QueryOneRecord(dest *interface{}, sql string, args ...interface{}) error {
	return c.Db.QueryRow(sql, args...).Scan(&dest)
}

// SQL Execute. example: create, insert, update, delete ...
//
// Arguments:
//	sql {string} - text for sql.
//	args {...interface{}} - args
//
// Returns:
//	{sql.Result} - results.
func (c *DatabaseOp) Execute(sql string, args ...interface{}) (sql.Result, error) {
	return c.Db.Exec(sql, args...)
}

// Create db table.
//
// warn: Columns must be defined programmatically due to the risk of SQL injection.
//
// Arguments:
//	tableName {string} - table name.
//	columns {string} - columns.
func (c *DatabaseOp) CreateTable(tableName string, columns string) error {
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS ? %s", columns)

	_, err := c.Execute(sql, tableName)
	if err != nil {
		return err
	}

	return nil
}

func (c *DatabaseOp) Close() {
	c.Db.Close()
}
