package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DatabaseOp struct {
	Db *sql.DB
}

// Create DB op.
//
// Arguments:
//	config {Config} - db config.
//
// Returns:
//	{*DatabaseOp} - database op instance.
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

// Count query.
//
// Arguments:
//	sql {string} - text for sql.
//	args {...interface{}} - args
//
// Returns:
//	{int} - count value.
func (c *DatabaseOp) Count(sql string, args ...interface{}) (int, error) {
	var count int
	if err := c.Db.QueryRow(sql, args...).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
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
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s %s", tableName, columns)

	_, err := c.Execute(sql)
	if err != nil {
		return err
	}

	return nil
}

func (c *DatabaseOp) Close() {
	c.Db.Close()
}
