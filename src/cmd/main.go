package main

import (
	"fmt"
	"log"
	"os"

	"github.com/JamesBLewis/go-pubsub-emulator/cmd/config"
	"github.com/JamesBLewis/go-pubsub-emulator/pkg/pubsubmanager"
)

func main() {
	// Read configuration from environment variables
	projectID := os.Getenv("PUBSUB_PROJECT")
	configFile := os.Getenv("CONFIG_FILE")

	if projectID == "" && configFile == "" {
		log.Fatal("Please provide either PROJECT_ID environment variable or CONFIG_FILE environment variable")
	}

	var topicConifg map[string][]string

	// If CONFIG_FILE is provided, read the configuration from the YAML file
	if configFile != "" {
		var err error
		topicConifg, err = config.ReadConfigFile(configFile)
		if err != nil {
			log.Fatalf("Error reading config file: %v", err)
		}
	} else if projectID != "" {
		// If PROJECT_ID is provided, handle command-line arguments (if any)
		topicID := os.Getenv("PUBSUB_TOPIC")
		subID := os.Getenv("PUBSUB_SUBSCRIPTION")

		if topicID == "" || subID == "" {
			log.Fatal("When providing PROJECT_ID, please also provide TOPIC_ID and SUB_ID")
		}

		// Add the specified topic and subscription to the map
		topicConifg = map[string][]string{
			topicID: {subID},
		}
	} else {
		log.Fatal("Invalid combination of environment variables. Please provide either CONFIG_FILE or PROJECT_ID along with TOPIC_ID and SUB_ID if needed.")
	}

	// Use the Pub/Sub manager to create topicConifg and subscriptions
	manager := pubsubmanager.NewPubSubManager(projectID)
	err := manager.CreateTopicsAndSubscriptions(topicConifg)
	if err != nil {
		log.Fatalf("Error creating topicConifg and subscriptions: %v", err)
	}

	fmt.Println("Topics and Subscriptions created successfully!")
}
