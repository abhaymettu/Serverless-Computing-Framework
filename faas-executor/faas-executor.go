package main

import (
	"log"
	"github.com/streadway/amqp"
)

func executeFunction(ch *amqp.Channel, msg amqp.Delivery) {
	// Extract function details from RabbitMQ message
	functionDetails := string(msg.Body)
	log.Printf("Function details extracted: %s", functionDetails)

	// Simulate using Docker SDK to run the function inside a container
	log.Printf("Executing function: %s", functionDetails)
	result := "Result of function execution" // placeholder

	// Send back the results to another queue
	err := ch.Publish(
		"",                   // exchange
		"faas-response-queue",// routing key
		false,                // mandatory
		false,                // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(result),
		})
	if err != nil {
		log.Printf("Failed to send a message to response queue: %v", err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@my-rabbitmq:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare the queue to ensure it exists
	_, err = ch.QueueDeclare(
		"faas-queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare the faas-queue: %v", err)
	}

	_, err = ch.QueueDeclare(
		"faas-response-queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare the faas-response-queue: %v", err)
	}

	msgs, err := ch.Consume(
		"faas-queue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			executeFunction(ch, msg)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
