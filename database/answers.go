package database

type Answers struct {
	Id     string `db:"id"`
	UserId string `db:"user_id"`
	Answer string `db:"answer"`
}

type DBAnswers struct {
	DB        *DatabaseOp
	TableName string
}

func NewDBAnswers(tableName string, config Config) (*DBAnswers, error) {
	db, err := NewDatabase(config)
	if err != nil {
		return nil, err
	}

	return &DBAnswers{
		DB:        db,
		TableName: tableName,
	}, nil
}

// Create table.
// If the table already exists, it will not be created.
func (c *DBAnswers) CreateTable() error {
	columns := "(id VARCHAR(256) NOT NULL, user_id VARCHAR(256) NOT NULL, answer VARCHAR(1024))"

	return c.DB.CreateTable(c.TableName, columns)
}

func (c *DBAnswers) AddAnswer(data Answers) error {
	sql := "INSERT INTO ? (id , user_id, answer) NOT NULL) VALUES (?, ?, ?)"

	_, err := c.DB.Execute(sql, c.TableName, data.Id, data.UserId, data.Answer)
	if err != nil {
		return err
	}
	return nil
}

func (c *DBAnswers) GetAnswers(targetId string) ([]Answers, error) {
	sql := "SELECT * FROM ? WHERE id = ?"
	result, err := c.DB.Query(sql, c.TableName, targetId)
	if err != nil {
		return nil, err
	}

	var answers []Answers = []Answers{}
	var data Answers

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

func (c *DBAnswers) Delete(targetId string) error {
	sql := "DELETE FROM ? WHERE id = '?'"

	_, err := c.DB.Execute(sql, c.TableName, targetId)
	if err != nil {
		return err
	}
	return nil
}
