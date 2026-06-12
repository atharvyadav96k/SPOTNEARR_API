package mq

import (
	"context"
	"log"
)

// Handler processes a raw JSON message body. Return non-nil to nack+requeue.
type Handler func(ctx context.Context, body []byte) error

// Subscriber consumes messages from the spotnearr.events exchange.
type Subscriber struct {
	conn *Conn
}

// NewSubscriber creates a Subscriber and verifies the exchange exists.
func NewSubscriber(c *Conn) (*Subscriber, error) {
	ch, err := c.raw.Channel()
	if err != nil {
		return nil, err
	}
	if err := declareExchange(ch); err != nil {
		ch.Close()
		return nil, err
	}
	ch.Close()
	return &Subscriber{conn: c}, nil
}

// Subscribe declares a durable queue, binds it to topic, and starts consuming
// in a background goroutine. The goroutine exits when ctx is cancelled.
func (s *Subscriber) Subscribe(ctx context.Context, queue string, topic Topic, handler Handler) error {
	ch, err := s.conn.raw.Channel()
	if err != nil {
		return err
	}

	// Fair dispatch: process one message at a time per consumer.
	if err := ch.Qos(1, 0, false); err != nil {
		ch.Close()
		return err
	}

	q, err := ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		ch.Close()
		return err
	}

	if err := ch.QueueBind(q.Name, string(topic), exchangeName, false, nil); err != nil {
		ch.Close()
		return err
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		ch.Close()
		return err
	}

	go func() {
		defer ch.Close()
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgs:
				if !ok {
					return
				}
				if err := handler(ctx, msg.Body); err != nil {
					log.Printf("mq: handler error on queue %s: %v", queue, err)
					msg.Nack(false, true) // requeue
				} else {
					msg.Ack(false)
				}
			}
		}
	}()

	return nil
}
