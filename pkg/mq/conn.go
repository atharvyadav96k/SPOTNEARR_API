package mq

import amqp "github.com/rabbitmq/amqp091-go"

const exchangeName = "spotnearr.events"

// Conn wraps an AMQP connection so callers never import amqp directly.
type Conn struct {
	raw *amqp.Connection
}

// Connect dials the RabbitMQ broker and returns a Conn.
func Connect(url string) (*Conn, error) {
	c, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	return &Conn{raw: c}, nil
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
