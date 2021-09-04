package database

type AbstractDBController struct {
	DB        *DatabaseOp
	TableName string
}

// Create table.
// If the table already exists, it will not be created.
func (c *AbstractDBController) CreateTable() error {
	columns := "(id VARCHAR(256) NOT NULL, user_id VARCHAR(256) NOT NULL)"

	return c.DB.CreateTable(c.TableName, columns)
}

// Clear all data in table.
func (c *AbstractDBController) ClearDB() error {
	sql := "DELETE FROM ?"

	_, err := c.DB.Execute(sql, c.TableName)
	if err != nil {
		return err
	}
	return nil
}
