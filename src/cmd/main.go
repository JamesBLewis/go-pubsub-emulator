package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/JamesBLewis/go-pubsub-emulator/cmd/config"
	"github.com/JamesBLewis/go-pubsub-emulator/pkg/healthcheck"
	"github.com/JamesBLewis/go-pubsub-emulator/pkg/pubsubmanager"
)

func main() {
	var cfg *config.Config
	ctx := context.Background()
	// load pubsub port from env
	port := os.Getenv("PUBSUB_PORT")
	if port == "" {
		log.Fatalln("Please provide PUBSUB_PORT environment variable")
	}

	project := os.Getenv("PUBSUB_PROJECT")
	if project == "" {
		log.Fatal("Please provide PUBSUB_PROJECT environment variable")
	}

	// for local testing
	var _ = os.Setenv("PUBSUB_EMULATOR_HOST", fmt.Sprintf("localhost:%s", port))

	// Attempt to read configuration from the file
	configFromFile, configErr := config.ReadConfigFile()
	if configErr != nil {
		log.Printf("error reading config file: %v", configErr)
	}
	if configFromFile != nil {
		// Configuration file found, use it
		log.Println("successfully read config from file")
		cfg = configFromFile
	} else {
		// Configuration file not found, check environment variables
		log.Println("unable to find config file, checking environment variables")

		topicID := os.Getenv("PUBSUB_TOPIC")
		subID := os.Getenv("PUBSUB_SUBSCRIPTION")

		if topicID == "" || subID == "" {
			log.Fatal("When providing PROJECT_ID, please also provide PUBSUB_TOPIC and PUBSUB_SUBSCRIPTION")
		}

		// Add the specified topic and subscription to the map
		cfg = &config.Config{
			Topics: map[string][]string{
				topicID: {subID},
			},
		}
	}

	// Perform Pub/Sub health check
	healthCheckErr := healthcheck.PollPubSub(port)
	if healthCheckErr != nil {
		log.Fatalf("Error performing Pub/Sub health check: %v", healthCheckErr)
	}

	// Use the Pub/Sub manager to create cfg and subscriptions
	manager, managerErr := pubsubmanager.NewPubSubManager(ctx, project)
	if managerErr != nil {
		log.Fatalf("Error creating Pub/Sub manager: %v", managerErr)
	}
	subscriptionErr := manager.CreateTopicsAndSubscriptions(ctx, cfg.Topics)
	if subscriptionErr != nil {
		log.Fatalf("Error creating cfg and subscriptions: %v", subscriptionErr)
	}

	fmt.Println("Topics and Subscriptions created successfully!")
}
