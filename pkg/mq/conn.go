package mq

import (
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const exchangeName = "spotnearr.events"

// Conn wraps an AMQP connection so callers never import amqp directly.
type Conn struct {
	raw *amqp.Connection
}

// Connect dials the RabbitMQ broker, retrying up to 10 times before giving up.
func Connect(url string) (*Conn, error) {
	const maxAttempts = 10
	const retryDelay = 3 * time.Second
	var err error
	for i := 1; i <= maxAttempts; i++ {
		var c *amqp.Connection
		c, err = amqp.Dial(url)
		if err == nil {
			log.Println("rabbitmq connection established")
			return &Conn{raw: c}, nil
		}
		log.Printf("rabbitmq not ready (attempt %d/%d): %v — retrying in %s", i, maxAttempts, err, retryDelay)
		time.Sleep(retryDelay)
	}
	return nil, err
}

// Close closes the underlying AMQP connection.
func (c *Conn) Close() error {
	return c.raw.Close()
}

// declareExchange ensures the shared topic exchange exists.
func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		exchangeName, // name
		"topic",      // kind
		true,         // durable
		false,        // auto-delete
		false,        // internal
		false,        // no-wait
		nil,
	)
}
