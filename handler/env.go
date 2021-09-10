package handler

import (
	"os"

	"github.com/hello-slide/synchronous-controller/database"
)

var dbUser = os.Getenv("DB_USER")
var dbPassword = os.Getenv("DB_PASSWORD")
var dbName = os.Getenv("DATABASE_NAME")

var dbConfig = database.NewConfig(dbUser, dbName, dbPassword)

var db *database.DatabaseOp

func Init() error {
	_db, err := database.NewDatabase(dbConfig)
	if err != nil {
		return err
	}
	db = _db

	return nil
}
