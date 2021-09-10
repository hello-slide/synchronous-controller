package database

import "fmt"

type Config struct {
	driverName     string
	dataSourceName string
}

// Create DB connect config.
//
// Arguments:
//	user {string} - User name.
//	dbName {string} - database name.
//	password {string} - password to connect database.
//
// Returns:
//	{*Config} - Database config.
func NewConfig(user string, dbName string, password string) *Config {
	config := fmt.Sprintf("host=project:region:instance user=%s dbname=%s password=%s sslmode=disable", user, dbName, password)

	return &Config{
		driverName:     "postgres",
		dataSourceName: config,
	}
}

// Create DB connect config by local.
//
// Arguments:
//	user {string} - User name.
//	dbName {string} - database name.
//	password {string} - password to connect database.
//
// Returns:
//	{*Config} - Database config.
func NewLocalConfig(user string, dbName string, password string) *Config {
	var config string

	if len(password) == 0 {
		config = fmt.Sprintf("user=%s dbname=%s sslmode=disable", user, dbName)
	} else {
		config = fmt.Sprintf("host=127.0.0.1 port=5432 user=%s dbname=%s password=%s sslmode=disable", user, dbName, password)
	}

	return &Config{
		driverName:     "postgres",
		dataSourceName: config,
	}
}
