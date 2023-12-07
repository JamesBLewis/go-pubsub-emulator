package main

import (
	"fmt"
	"log"
	"os"

	"github.com/JamesBLewis/go-pubsub-emulator/cmd/config"
	"github.com/JamesBLewis/go-pubsub-emulator/pkg/pubsubmanager"
)

func main() {
	var cfg *config.Config
	// Attempt to read configuration from the file
	configFromFile, err := config.ReadConfigFile()
	if err == nil {
		// Configuration file found, use it
		log.Println("successfully read config from file")
		cfg = configFromFile
	} else {
		// Configuration file not found, check environment variables
		log.Println("unable to read config file, checking environment variables")
		project := os.Getenv("PUBSUB_PROJECT")
		if project == "" {
			log.Fatal("Please provide PUBSUB_PROJECT environment variable")
		}

		topicID := os.Getenv("PUBSUB_TOPIC")
		subID := os.Getenv("PUBSUB_SUBSCRIPTION")

		if topicID == "" || subID == "" {
			log.Fatal("When providing PROJECT_ID, please also provide PUBSUB_TOPIC and PUBSUB_SUBSCRIPTION")
		}

		// Add the specified topic and subscription to the map
		cfg = &config.Config{
			ProjectID: project,
			Topics: map[string][]string{
				topicID: {subID},
			},
		}
	}

	// Use the Pub/Sub manager to create cfg and subscriptions
	manager := pubsubmanager.NewPubSubManager(cfg.ProjectID)
	err = manager.CreateTopicsAndSubscriptions(cfg.Topics)
	if err != nil {
		log.Fatalf("Error creating cfg and subscriptions: %v", err)
	}

	fmt.Println("Topics and Subscriptions created successfully!")
}
