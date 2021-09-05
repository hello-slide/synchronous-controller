package pubsub_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	_pubsub "cloud.google.com/go/pubsub"
	"github.com/hello-slide/synchronous-controller/pubsub"
)

func TestPubSub(t *testing.T) {
	if os.Getenv("LOCAL_TEST") != "pubsub" {
		return
	}

	ctx := context.Background()

	keyPath := os.Getenv("IAM_PATH")
	projectId := "helloslide"
	topicName := "test"
	subId := "subb"

	client, err := pubsub.CreateClientLocal(ctx, projectId, keyPath)
	if err != nil {
		t.Fatalf("create client error: %v", err)
	}

	pubsubOpPublish := pubsub.NewPubSub(ctx, client)
	pubsubOpPublish, err = pubsubOpPublish.CreateTopic(topicName)
	if err != nil {
		t.Fatalf("create topic error: %v", err)
	}

	time.Sleep(1 * time.Minute)

	for index, message := range []string{"hogehoge", "hugahuga"} {
		done := make(chan bool)
		go publish(client, t, topicName, message, done)

		// subscription
		pubsubOp, err := pubsub.NewPubSub(ctx, client).SetTopic(topicName)
		if err != nil {
			t.Fatalf("create topic error: %v", err)
		}

		exist, err := pubsubOp.ExistTopic()
		if err != nil {
			t.Fatalf("exist error: %v", err)
		}
		if !exist {
			t.Fatal("topic dose not exists.")
		}

		var resMessage string

		if index == 0 {
			resMessage, err = pubsubOp.CreateSubscription(subId)
		} else {
			resMessage, err = pubsubOp.Subscription(subId)
		}

		<-done

		if err != nil {
			t.Fatalf("subscription error: %v", err)
		}

		if resMessage != message {
			t.Fatalf("It is different from the published message.")
		}
	}
}

func publish(client *_pubsub.Client, t *testing.T, topicName string, message string, done chan bool) {
	time.Sleep(10 * time.Second)
	fmt.Println("publish!")

	ctx := context.Background()

	pubsubOpPublish := pubsub.NewPubSub(ctx, client)

	var err error
	pubsubOpPublish, err = pubsubOpPublish.SetTopic(topicName)
	if err != nil {
		t.Fatalf("create topic error: %v", err)
	}

	id, err := pubsubOpPublish.Publish([]byte(message))
	if err != nil {
		t.Fatalf("publish error: %v", err)
	}
	t.Logf("success publish! id: %v", id)
	done <- true
}
