package pubsub_test

import (
	"context"
	"os"
	"testing"
	"time"

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

	client, err := pubsub.CreateClientLocal(ctx, projectId, keyPath)
	if err != nil {
		t.Fatalf("create client error: %v", err)
	}

	for _, message := range []string{"hogehoge", "hooaaaa"} {
		// publish
		pubsubOp, err := pubsub.NewPubSub(ctx, client).CreateTopic(topicName)
		if err != nil {
			t.Fatalf("create topic error: %v", err)
		}

		id, err := pubsubOp.Publish([]byte(message))
		if err != nil {
			t.Fatalf("publish error: %v", err)
		}
		t.Logf("success publish! id: %v", id)

		// subscription
		pubsubOp, err = pubsub.NewPubSub(ctx, client).SetTopic(topicName)
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

		extension := 30 * time.Second
		resMessage, err := pubsubOp.Subscription(extension)
		if err != nil {
			t.Fatalf("subscription error: %v", err)
		}

		if resMessage != message {
			t.Fatalf("It is different from the published message.")
		}
	}
}
