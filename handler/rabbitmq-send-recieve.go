package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SendDataToMQ() {
	conn, err := amqp.Dial("amqp://user:1234@localhost:5672/")
	if err != nil {
		fmt.Println("connecetion failed ", err)
	} else {
		fmt.Println("connection to rabbit mq established")
	}
	defer conn.Close()

	channel, _ := conn.Channel()
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println("failed to declare queue", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	finalJSON, _ := json.Marshal(bulkOrders)

	err = channel.PublishWithContext(
		ctx,
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        finalJSON,
		},
	)
	if err != nil {
		fmt.Println("message not published", err)
	}

	log.Printf("Message Published: %s\n", string(finalJSON))
}

func ReceiveFromMQ() {
	conn, err := amqp.Dial("amqp://user:1234@localhost:5672/")
	if err != nil {
		fmt.Println("connection failed", err)
	}
	defer conn.Close()

	channel, _ := conn.Channel()
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println("failed to declare queue", err)
	}

	msgs, err := channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println("consumer not created: ", err)
	}

	var block = make(chan struct{})

	go func() {
		for msg := range msgs {
			log.Printf("Received a message: %s", msg.Body)
			processMessage(msg.Body)
		}
	}()

	<-block
}
