package main

import (
	"fmt"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var messages []string // In-memory store for messages

func main() {
	// Kafka consumer configuration
	config := &kafka.ConfigMap{
		"bootstrap.servers": "kafka:9092", // Replace with your Kafka broker(s)
		"group.id":          "my-group",   // Consumer group ID
		"auto.offset.reset": "earliest",   // Start reading from the earliest offset
	}

	// Create a new Kafka consumer
	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		panic(fmt.Sprintf("Failed to create consumer: %v", err))
	}
	defer consumer.Close()

	// Subscribe to the Kafka topic
	topic := "my-topic"
	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		panic(fmt.Sprintf("Failed to subscribe to topic %s: %v", topic, err))
	}

	// Goroutine to consume messages
	go func() {
		for {
			msg, err := consumer.ReadMessage(-1) // Block until a message is received
			if err != nil {
				fmt.Printf("Consumer error: %v\n", err)
				continue
			}
			fmt.Printf("Received message: %s\n", string(msg.Value))
			messages = append(messages, string(msg.Value)) // Store the message
		}
	}()

	// HTTP handler to display messages
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<h1>Kafka Messages</h1>")
		for _, msg := range messages {
			fmt.Fprintf(w, "<p>%s</p>", msg)
		}
	})

	// Start the HTTP server
	fmt.Println("Starting HTTP server on :9090")
	if err := http.ListenAndServe(":9090", nil); err != nil {
		panic(fmt.Sprintf("Failed to start HTTP server: %v", err))
	}
}
