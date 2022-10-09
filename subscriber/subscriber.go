package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

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

	prefetchCount, err := strconv.Atoi(os.Args[1])
	if err != nil {
		failOnError(err, "Failed to get lenOfBytes")
	}

	ch.Qos(prefetchCount)

	err = ch.ExchangeDeclare(
		"bytes",  // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"bytes", // name
		false,   // durable
		false,   // delete when unused
		true,    // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,  // queue name
		"",      // routing key
		"bytes", // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	csvFile, err := os.Create("execs.csv")
	defer csvFile.Close()
	csvwriter := csv.NewWriter(csvFile)
	csvwriter.Write([]string{"Bytes", "ExecTime"})
	for msg := range msgs {
		msg.Ack(false)
		arrivalTime := time.Now().UnixNano()
		message := types.Message{}
		err := json.Unmarshal(msg.Body, &message)
		if err != nil {
			failOnError(err, "Failed to receive a message")
		}

		csvwriter.Write([]string{fmt.Sprintf("%d", len(message.Content)), fmt.Sprintf("%d", message.CreatedAt-arrivalTime)})
	}
}
