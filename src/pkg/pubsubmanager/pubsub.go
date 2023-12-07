package pubsubmanager

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/pubsub"
)

type PubSubManager struct {
	projectID string
	client    *pubsub.Client
}

func NewPubSubManager(ctx context.Context, projectID string) (*PubSubManager, error) {
	// Create a Pub/Sub client with the specified project ID
	pubsubClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("error creating Pub/Sub client: %v", err)
	}

	return &PubSubManager{
		projectID: projectID,
		client:    pubsubClient,
	}, nil
}

func (m *PubSubManager) createTopic(ctx context.Context, topicID string) (*pubsub.Topic, error) {
	topic := m.client.Topic(topicID)
	exists, err := topic.Exists(ctx)
	if err != nil {
		return nil, fmt.Errorf("error checking if topic exists: %v", err)
	}

	if !exists {
		topic, err = m.client.CreateTopic(ctx, topicID)
		if err != nil {
			return nil, fmt.Errorf("error creating topic: %v", err)
		}
		log.Printf("Topic %s created", topicID)
	} else {
		log.Printf("Topic %s already exists", topicID)
	}

	return topic, nil
}

func (m *PubSubManager) createSubscription(ctx context.Context, topic *pubsub.Topic, subID string) (*pubsub.Subscription, error) {
	sub := m.client.Subscription(subID)
	exists, err := sub.Exists(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error checking if subscription exists: %v", err)
	}

	if !exists {
		sub, err = m.client.CreateSubscription(ctx, subID, pubsub.SubscriptionConfig{
			Topic:       topic,
			AckDeadline: 10 * time.Second, // GCP default Ack deadline
		})
		if err != nil {
			return nil, fmt.Errorf("Error creating subscription: %v", err)
		}
		log.Printf("Subscription %s created", subID)
	} else {
		log.Printf("Subscription %s already exists", subID)
	}

	return sub, nil
}

func (m *PubSubManager) CreateTopicsAndSubscriptions(ctx context.Context, topics map[string][]string) error {
	// Create topics and subscriptions
	for topicID, subs := range topics {
		topic, err := m.createTopic(ctx, topicID)
		if err != nil {
			return fmt.Errorf("error creating topic: %v", err)
		}

		for _, subID := range subs {
			_, err := m.createSubscription(ctx, topic, subID)
			if err != nil {
				return fmt.Errorf("error creating subscription: %v", err)
			}
		}
	}

	return nil
}
