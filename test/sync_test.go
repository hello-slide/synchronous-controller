package test

import (
	"testing"

	"github.com/hello-slide/synchronous-controller/host"
	"github.com/hello-slide/synchronous-controller/visitor"
)

func TestSync(t *testing.T) {
	ip := "192.168.0.1"

	// create session from host.
	token := host.CreateSession(ip)

	// Get topic from visitor.
	data, err := visitor.GetTopic(token)
	if err != nil {
		t.Fatalf("Get topic error: %v", err)
	}
	if data != "" {
		t.Fatalf("topic already exists. data is %v", data)
	}

	// Set topics from host.
	topic := "hogehoge"
	if err := host.SetTopic(token, topic); err != nil {
		t.Fatalf("Set topic error: %v", err)
	}

	// Get topic from visitor.
	data, err = visitor.GetTopic(token)
	if err != nil {
		t.Fatalf("Get topic error: %v", err)
	}
	if data != "hogehoge" {
		t.Fatalf("The topic is different. data is %v", data)
	}

	// Set answer from visitor.
	ansOne := "1"
	ansTwo := "2"
	if err := visitor.AddAnswer(token, ansOne); err != nil {
		t.Fatalf("Add answer error: %v", err)
	}
	if err := visitor.AddAnswer(token, ansTwo); err != nil {
		t.Fatalf("Add answer error: %v", err)
	}

	// Get answer from host.
	answers, err := host.GetResult(token)
	if err != nil {
		t.Fatalf("Get answer error: %v", err)
	}

	if len(answers) != 2 {
		t.Fatalf("The number of answers is different. length is %v. in %v", len(answers), answers)
	}
	for _, element := range answers {
		switch element {
		case ansOne:
			break
		case ansTwo:
			break
		default:
			t.Fatalf("There is an answer I don't know. %v", element)
		}
	}

	// Close session from host.
	if err := host.Close(token); err != nil {
		t.Fatalf("Session close error: %v", err)
	}

	// Check if the session does not exist from visitor.
	_, err = visitor.GetTopic(token)

	if err == nil {
		t.Fatalf("The session should not exist, but it does.")
	}
}
