package mq

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Publisher sends messages to the spotnearr.events exchange.
type Publisher struct {
	conn *Conn
}

// NewPublisher creates a Publisher and verifies the exchange exists.
func NewPublisher(c *Conn) (*Publisher, error) {
	ch, err := c.raw.Channel()
	if err != nil {
		return nil, err
	}
	if err := declareExchange(ch); err != nil {
		ch.Close()
		return nil, err
	}
	ch.Close()
	log.Println("rabbitmq publisher ready")
	return &Publisher{conn: c}, nil
}

// Publish JSON-encodes payload and routes it by topic to the exchange.
// A new channel is opened per call — safe for concurrent use.
func (p *Publisher) Publish(ctx context.Context, topic Topic, payload any) error {
	ch, err := p.conn.raw.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return ch.PublishWithContext(ctx,
		exchangeName,  // exchange
		string(topic), // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
}
