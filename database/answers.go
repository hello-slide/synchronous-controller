package database

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

	return &DBAnswers{
		AbstractDBController{
			DB:        db,
			TableName: tableName,
		},
	}, nil
}

// Insert the new answer.
//
// Arguments:
// data {*Answer} - Answer data.
func (c *DBAnswers) AddAnswer(data *Answer) error {
	sql := "INSERT INTO ? (id , user_id, answer) NOT NULL) VALUES (?, ?, ?)"

	_, err := c.DB.Execute(sql, c.TableName, data.Id, data.UserId, data.Answer)
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
	sql := "SELECT * FROM ? WHERE id = ?"
	result, err := c.DB.Query(sql, c.TableName, targetId)
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
	sql := "DELETE FROM ? WHERE id = '?'"

	_, err := c.DB.Execute(sql, c.TableName, targetId)
	if err != nil {
		return err
	}
	return nil
}
