package handler

import (
	"os"

	"github.com/hello-slide/synchronous-controller/database"
)

var dbUser = os.Getenv("DB_USER")
var dbPassword = os.Getenv("DB_PASSWORD")
var dbName = os.Getenv("DATABASE_NAME")

var dbConfig = database.NewConfig(dbUser, dbName, dbPassword)
