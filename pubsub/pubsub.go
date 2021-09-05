package pubsub

import (
	"context"
	"errors"
	"os"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

// Create pubsub client.
func CreateClient(ctx context.Context) (*pubsub.Client, error) {
	key := []byte(os.Getenv("KEY"))
	projectId := os.Getenv("PUBSUB_PROJECT_ID")

	return pubsub.NewClient(ctx, projectId, option.WithCredentialsJSON(key))
}

// Create pubsub used by local.
func CreateClientLocal(ctx context.Context, projectId string, keyPath string) (*pubsub.Client, error) {
	return pubsub.NewClient(ctx, projectId, option.WithCredentialsFile(keyPath))
}

type PubSubController struct {
	ctx    context.Context
	client *pubsub.Client
	topic  *pubsub.Topic
}

// Create pubsub controller instance.
//
// Arguments:
//	ctx {context.Context} - context.
//	client {*pubsub.Client} - client of pubsub.
//
// Returns:
//	{*PubSubController} - pubsub controller instance.
func NewPubSub(ctx context.Context, client *pubsub.Client) *PubSubController {
	return &PubSubController{
		ctx:    ctx,
		client: client,
		topic:  nil,
	}
}

// Create a new topic.
//
// Arguments:
//	topicName {string} - topic name.
//
// Returns:
//	{*PubSubController} - pubsub controller instance.
func (c *PubSubController) CreateTopic(topicName string) (*PubSubController, error) {
	topic, err := c.client.CreateTopic(c.ctx, topicName)
	if err != nil {
		return nil, err
	}

	return &PubSubController{
		ctx:    c.ctx,
		client: c.client,
		topic:  topic,
	}, nil
}

// Set already exist topic.
//
// Arguments:
//	topicName {string} - topic name.
//
// Returns:
//	{*PubSubController} - pubsub controller instance.
func (c *PubSubController) SetTopic(topicName string) (*PubSubController, error) {
	topic := c.client.Topic(topicName)

	exist, err := topic.Exists(c.ctx)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New("Topic dose not exist.")
	}

	return &PubSubController{
		ctx:    c.ctx,
		client: c.client,
		topic:  topic,
	}, nil
}

// Check the existence of the topic.
func (c *PubSubController) checkTopic() error {
	if c.topic == nil {
		return errors.New("the topic is unknown")
	}
	return nil
}

// Stop topic.
func (c *PubSubController) StopTopic() error {
	if err := c.checkTopic(); err != nil {
		return err
	}

	c.topic.Stop()
	return nil
}

// Publish to pubsub.
//
// Arguments:
//	message {[]byte} - send message.
//
// Returns:
//	{string} - message id.
func (c *PubSubController) Publish(message []byte) (string, error) {
	if err := c.checkTopic(); err != nil {
		return "", err
	}

	res := c.topic.Publish(c.ctx, &pubsub.Message{
		Data: message,
	})

	return res.Get(c.ctx)
}

// check if exist topics.
//
// Returns:
//	{bool} - exist if true, false is not.
func (c *PubSubController) ExistTopic() (bool, error) {
	if err := c.checkTopic(); err != nil {
		return false, err
	}

	return c.topic.Exists(c.ctx)
}

// Create a new subscription and subscribe.
//
// Arguments:
//	id {string} - subscription id.
//	extension {time.Duration} - Time limit
//
// Returns:
//	{string} - received message.
func (c *PubSubController) CreateSubscription(id string) (string, error) {
	if err := c.checkTopic(); err != nil {
		return "", err
	}

	sub, err := c.client.CreateSubscription(c.ctx, id, pubsub.SubscriptionConfig{
		Topic: c.topic,
	})
	if err != nil {
		return "", err
	}

	return c.revive(sub)
}

// Subscription by already exist sub id.
// Execution will be blocked until you get a new topic or the time limit is reached.
//
// Arguments:
//	id {string} - subscription id.
//	extension {time.Duration} - Time limit
//
// Returns:
//	{string} - received message.
func (c *PubSubController) Subscription(id string) (string, error) {
	sub := c.client.Subscription(id)

	return c.revive(sub)
}

func (c *PubSubController) revive(sub *pubsub.Subscription) (string, error) {
	var message string

	sub.ReceiveSettings.Synchronous = true

	err := sub.Receive(c.ctx, func(ctx context.Context, m *pubsub.Message) {
		message = string(m.Data)
		m.Ack()
	})

	if err != context.Canceled {
		return "", err
	}

	return message, err
}
