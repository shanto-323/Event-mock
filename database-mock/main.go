package main

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/tinrab/retry"
)

type Config struct {
	RabbitmqUrl string `envconfig:"RABBITMQ_URL"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	var conn *amqp.Connection
	retry.ForeverSleep(
		2*time.Second,
		func(_ int) error {
			conn, err = amqp.Dial(cfg.RabbitmqUrl)
			if err != nil {
				return err
			}
			return nil
		},
	)
	defer conn.Close()
	repo := NewDbRepository()

	event, err := NewEventPublisher(conn, repo)
	if err != nil {
		log.Fatal(err)
	}

	go event.Consume(GATEWAY_EVENT_GET)
	go event.Consume(GATEWAY_EVENT_CREATE)
	select {}
}
