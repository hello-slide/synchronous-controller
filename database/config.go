package database

import "fmt"

type Config struct {
	driverName     string
	dataSourceName string
}

func NewConfig(user string, dbName string, password string) *Config {
	config := fmt.Sprintf("host=project:region:instance user=%s dbname=%s password=%s sslmode=disable", user, dbName, password)

	return &Config{
		driverName:     "cloudsqlpostgres",
		dataSourceName: config,
	}
}

func NewLocalConfig(user string, dbName string, password string) *Config {
	config := fmt.Sprintf("host=localhost user=%s dbname=%s password=%s sslmode=disable", user, dbName, password)

	return &Config{
		driverName:     "postgres",
		dataSourceName: config,
	}
}
