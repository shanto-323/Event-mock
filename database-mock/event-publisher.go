package main

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	GATEWAY_EVENT_GET    = "gate.get"
	GATEWAY_EVENT_CREATE = "gate.create"
	NOTIFICATION_CREATE  = "notification.create"
)

type EventPublisher struct {
	ch   *amqp.Channel
	repo Repository
}

func NewEventPublisher(conn *amqp.Connection, repo Repository) (*EventPublisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &EventPublisher{
		ch:   ch,
		repo: repo,
	}, nil
}

func (e *EventPublisher) Close() error {
	return e.ch.
		Close()
}

func (e *EventPublisher) Produce(queueType string, payload []byte) error {
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
			e.handleMessage(queueType, m)
		}()
	}
	<-forever
}

func (e *EventPublisher) handleMessage(queueName string, d amqp.Delivery) {
	switch queueName {
	case GATEWAY_EVENT_CREATE:
		{
			resp := &CreateModel{}
			err := json.Unmarshal(d.Body, resp)
			if err != nil {
				log.Panic(err)
			}
			e.repo.MockCreateData(resp)
			go e.Produce(NOTIFICATION_CREATE, []byte("data stored in database"))
		}

	case GATEWAY_EVENT_GET:
		{
			resp := &GetModel{}
			err := json.Unmarshal(d.Body, resp)
			if err != nil {
				log.Panic(err)
			}
			e.repo.MockGetData(resp)
			go e.Produce(NOTIFICATION_CREATE, []byte("got data from database"))
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
