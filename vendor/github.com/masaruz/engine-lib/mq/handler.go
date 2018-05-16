package mq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

// MQ to handle message queue
type MQ struct {
	q    amqp.Queue
	ch   *amqp.Channel
	conn *amqp.Connection
}

// Init connection
func (m *MQ) Init(ip string, port int) {
	var err error
	m.conn, err = amqp.Dial(fmt.Sprintf("amqp://guest:guest@%s:%d/", ip, port))
	failOnError(err, "Failed to connect to RabbitMQ")
	m.ch, err = m.conn.Channel()
	failOnError(err, "Failed to open a channel")
}

// Create initiate queue
func (m *MQ) Create(queue string) {
	var err error
	m.q, err = m.ch.QueueDeclare(
		queue, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")
}

// Send the message to queue
func (m *MQ) Send() {
	body := "hello"
	err := m.ch.Publish(
		"",       // exchange
		m.q.Name, // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
}

// Receive message from queue
func (m *MQ) Receive() {
	msgs, err := m.ch.Consume(
		m.q.Name, // queue
		"",       // consumer
		true,     // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)
	failOnError(err, "Failed to register a consumer")

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)
	}
}

// Destroy close everything
func (m *MQ) Destroy() {
	m.ch.Close()
	m.conn.Close()
}
