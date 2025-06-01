package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	GATEWAY_EVENT_GET    = "gate.get"
	GATEWAY_EVENT_CREATE = "gate.create"
)

type EventPublisher struct {
	ch *amqp.Channel
}

func NewEventPublisher(conn *amqp.Connection) (*EventPublisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &EventPublisher{
		ch: ch,
	}, nil
}

func (e *EventPublisher) Close() error {
	return e.ch.
		Close()
}

func (e *EventPublisher) Producer(queueType string, payload []byte) error {
	ch := e.ch
	queue, err := getQueue(queueType, ch)
	if err != nil {
		return err
	}

	return ch.Publish(
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payload,
		},
	)
}

func getQueue(queueType string, ch *amqp.Channel) (amqp.Queue, error) {
	queue, err := ch.QueueDeclare(
		queueType,
		false,
		false,
		false,
		false,
		amqp.Table{
			"x-delivery-limit": 2, // Max 2 attempts (RabbitMQ 3.8+)
		},
	)
	return queue, err
}
