package main

import (
	"context"
	"log"
	"os"
	"time"
	"strconv"
	"encoding/json"

	"github.com/jgfn1/golang-rabbitmq/types"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"bytes",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	lenOfBytes, err := strconv.Atoi(os.Args[1])
	if err != nil {
		failOnError(err, "Failed to get lenOfBytes")
	}

	for i := 0; i < 10000; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		
		body := bodyFrom(os.Args, lenOfBytes)
		err = ch.PublishWithContext(ctx,
			"bytes", // exchange
			"",     // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body: body,
			})
		failOnError(err, "Failed to publish a message")
	}
}

func bodyFrom(args []string, lenOfBytes int) ([]byte) {
	
	message, err := buildMessage(lenOfBytes)
	if err != nil {
		failOnError(err, "Failed to create message")
	}

	return message
}

func buildMessage(lenOfBytes int) ([]byte, error) {
	return json.Marshal(types.Message{
		Content: make([]byte, lenOfBytes),
		CreatedAt: time.Now().UnixNano(),
	})
}
