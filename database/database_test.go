package database_test

import (
	"os"
	"testing"

	"github.com/hello-slide/synchronous-controller/database"
)

func TestConnectUser(t *testing.T) {
	if len(os.Getenv("LOCAL_TEST")) == 0 {
		return
	}

	// ----config------

	user := os.Getenv("DB_USER")
	password := ""
	dbName := "hello-slide-test"
	tableName := "connectusers"

	rootId := "rootIdsample"

	sampleUserIdOne := "samplesample1"
	sampleUserIdTwo := "samplesample2"

	sampleDataOne := &database.ConnectUser{
		Id:     rootId,
		UserId: sampleUserIdOne,
	}

	sampleDataTwo := &database.ConnectUser{
		Id:     rootId,
		UserId: sampleUserIdTwo,
	}

	// -----end------

	config := database.NewLocalConfig(user, dbName, password)

	connectUser, err := database.NewDBConnectUsers(tableName, *config)
	if err != nil {
		t.Fatalf("db connect error: %v", err)
	}

	if err := connectUser.CreateTable(); err != nil {
		t.Fatalf("create table error: %v", err)
	}

	// Add users.
	if err := connectUser.AddUser(*sampleDataOne); err != nil {
		t.Fatalf("add user error: %v", err)
	}
	if err := connectUser.AddUser(*sampleDataTwo); err != nil {
		t.Fatalf("add user second error: %v", err)
	}

	users, err := connectUser.GetUserNumber(rootId)
	if err != nil {
		t.Fatalf("get users error: %v", err)
	}

	if users != 2 {
		t.Fatalf("The number of users stored in the database is different. num: %v", users)
	}

	if err := connectUser.Delete(rootId); err != nil {
		t.Fatalf("delete error: %v", err)
	}
}
