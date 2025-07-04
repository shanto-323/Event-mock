package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	GATEWAY_EVENT_GET    = "gate.get"
	GATEWAY_EVENT_CREATE = "gate.create"
	NOTIFICATION_CREATE  = "notification.create"
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

func (e *EventPublisher) Consume(queueType string) {
	ch := e.ch
	queue, err := getQueue(queueType, ch)
	if err != nil {
		log.Print(err.Error())
	}

	msg, err := ch.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Print(err.Error())
	}
	forever := make(chan interface{})
	for m := range msg {
		go func() {
			e.handleMessage(queue.Name, m)
		}()
	}
	<-forever
}

func (e *EventPublisher) handleMessage(queueName string, d amqp.Delivery) {
	switch queueName {
	case GATEWAY_EVENT_CREATE:
		{
			log.Println("sent data for creation", string(d.Body))
		}

	case GATEWAY_EVENT_GET:
		{
			log.Println("request for getting data", string(d.Body))
		}
	case NOTIFICATION_CREATE:
		{
			log.Println("new notification", string(d.Body))
		}
	}
}

func getQueue(queueType string, ch *amqp.Channel) (amqp.Queue, error) {
	queue, err := ch.QueueDeclare(
		queueType,
		false,
		false,
		false,
		false,
		nil,
	)
	return queue, err
}
