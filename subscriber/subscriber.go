package main

import (
	"fmt"
	"time"
	"log"
	"os"
	"encoding/json"
	"encoding/csv"

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

	q, err := ch.QueueDeclare(
		"bytes",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		"bytes", // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	csvFile, err := os.Create("execs.csv")
	csvwriter := csv.NewWriter(csvFile)
	csvwriter.Write([]string{"Bytes", "ExecTime"})
	for msg := range msgs {
		message := types.Message{}
		err := json.Unmarshal(msg.Body, &message)
		if err != nil {
			failOnError(err, "Failed to receive a message")
		}

		csvwriter.Write([]string{fmt.Sprintf("%d", len(message.Content)), fmt.Sprintf("%d", (time.Now().UnixNano() - message.CreatedAt))})
	}
	csvFile.Close()
}
