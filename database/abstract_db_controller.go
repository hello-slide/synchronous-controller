package database

import "fmt"

type AbstractDBController struct {
	DB        *DatabaseOp
	TableName string
	Columns   string
}

// Create table.
// If the table already exists, it will not be created.
func (c *AbstractDBController) CreateTable() error {
	return c.DB.CreateTable(c.TableName, c.Columns)
}

// Clear all data in table.
func (c *AbstractDBController) ClearDB() error {
	sql := fmt.Sprintf("DELETE FROM %s", c.TableName)

	_, err := c.DB.Execute(sql, c.TableName)
	if err != nil {
		return err
	}
	return nil
}
