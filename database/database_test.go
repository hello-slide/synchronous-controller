package database_test

import (
	"log"
	"os"
	"testing"

	"github.com/hello-slide/synchronous-controller/database"
)

var db *database.DatabaseOp

func init() {
	if os.Getenv("LOCAL_TEST") != "db" {
		log.Printf("Database test was not run.")
		return
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PW")
	dbName := "hello-slide-test"

	log.Printf("username: %v", user)
	log.Printf("password: %v", password)

	config := database.NewLocalConfig(user, dbName, password)
	_db, err := database.NewDatabase(config)
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}
	db = _db
}

func TestConnectUser(t *testing.T) {
	if os.Getenv("LOCAL_TEST") != "db" || db == nil {
		t.Log("not run test")
		return
	}
	t.Log("run test")

	// ----config------
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

	connectUser := database.NewDBConnectUsers(tableName, db)

	if err := connectUser.CreateTable(); err != nil {
		t.Fatalf("create table error: %v", err)
	}

	// Add users.
	if err := connectUser.AddUser(sampleDataOne); err != nil {
		t.Fatalf("add user error: %v", err)
	}
	if err := connectUser.AddUser(sampleDataTwo); err != nil {
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

func TestAnswers(t *testing.T) {
	if os.Getenv("LOCAL_TEST") != "db" || db == nil {
		t.Log("not run test")
		return
	}
	t.Log("run test")

	// ----config------
	tableName := "answers"

	rootId := "rootIdsample"

	sampleDataOne := &database.Answer{
		Id:     rootId,
		UserId: "samplesample1",
		Answer: "hogehoge",
	}

	sampleDataTwo := &database.Answer{
		Id:     rootId,
		UserId: "samplesample2",
		Answer: "hugahuga",
	}

	// -----end------

	dbAnswer := database.NewDBAnswers(tableName, db)

	if err := dbAnswer.CreateTable(); err != nil {
		t.Fatalf("create table error: %v", err)
	}

	if err := dbAnswer.AddAnswer(sampleDataOne); err != nil {
		t.Fatalf("add answers error: %v", err)
	}

	if err := dbAnswer.AddAnswer(sampleDataTwo); err != nil {
		t.Fatalf("add answers error: %v", err)
	}

	answers, err := dbAnswer.GetAnswers(rootId)
	if err != nil {
		t.Fatalf("get answer error: %v", err)
	}

	if len(answers) != 2 {
		t.Fatalf("The number of answers stored in the database is different. num: %v", len(answers))
	}

	if err := dbAnswer.Delete(rootId); err != nil {
		t.Fatalf("delete error: %v", err)
	}
}

func TestTopic(t *testing.T) {
	if os.Getenv("LOCAL_TEST") != "db" || db == nil {
		t.Log("not run test")
		return
	}
	t.Log("run test")

	// ----config------
	tableName := "topic"

	rootId := "topic_id"

	sampleData := &database.Topic{
		Id:       rootId,
		Topic:    "hoge",
		IsUpdate: false,
	}

	sampleData2 := &database.Topic{
		Id:       rootId,
		Topic:    "huga",
		IsUpdate: true,
	}

	// -----end------

	dbTopic := database.NewDBTopic(tableName, db)

	if err := dbTopic.CreateTable(); err != nil {
		t.Fatalf("create table error: %v", err)
	}

	if err := dbTopic.CreateTopic(sampleData); err != nil {
		t.Fatalf("create topic error: %v", err)
	}

	updateFlag, err := dbTopic.GetIsUpdate(rootId)
	if err != nil {
		t.Fatalf("get update flag error: %v", err)
	}
	if updateFlag {
		t.Fatal("error: updateflag is not false.")
	}

	result, err := dbTopic.GetTopic(rootId)
	if err != nil {
		t.Fatalf("get topic error: %v", err)
	}
	if result != "hoge" {
		t.Fatalf("the result is different. value: %v", result)
	}

	// update topic
	if err := dbTopic.UpdateTopic(sampleData2); err != nil {
		t.Fatalf("update topic error: %v", err)
	}

	updateFlag, err = dbTopic.GetIsUpdate(rootId)
	if err != nil {
		t.Fatalf("get update flag error: %v", err)
	}
	if !updateFlag {
		t.Fatal("error: updateflag is true.")
	}

	result, err = dbTopic.GetTopic(rootId)
	if err != nil {
		t.Fatalf("get topic error: %v", err)
	}
	if result != "huga" {
		t.Fatalf("the result is different. value: %v", result)
	}

	if err := dbTopic.Delete(rootId); err != nil {
		t.Fatalf("delete error: %v", err)
	}

	_, err = dbTopic.GetIsUpdate(rootId)
	if err == nil {
		t.Fatal("The topic id has been removed and should be an error.")
	}
}
