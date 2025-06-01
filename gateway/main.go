package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kelseyhightower/envconfig"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/tinrab/retry"
)

type Config struct {
	RabbitmqUrl string `envconfig:"RABBITMQ_URL"`
}

type Server struct {
	Event  *EventPublisher
	IpAddr string
}

func NewServer(event *EventPublisher, ipAddr int) *Server {
	ip := fmt.Sprintf(":%d", ipAddr)
	return &Server{
		Event:  event,
		IpAddr: ip,
	}
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
	event, err := NewEventPublisher(conn)
	if err != nil {
		log.Fatal(err)
	}

	server := NewServer(event, 8080)

	router := chi.NewMux()
	router.HandleFunc("/create", handleFunc(server.HandleCreateUser))
	router.HandleFunc("/get", handleFunc(server.HandleGetUser))
	log.Fatal(http.ListenAndServe(server.IpAddr, router))
}

func (s *Server) HandleCreateUser(w http.ResponseWriter, r *http.Request) error {
	payload, err := json.Marshal(
		CreateModel{
			ID:     "123456",
			Name:   "Bot1",
			Amount: 100000,
		},
	)
	if err != nil {
		return err
	}

	if err := s.Event.Producer(
		GATEWAY_EVENT_CREATE,
		payload,
	); err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, "in process create")
}

func (s *Server) HandleGetUser(w http.ResponseWriter, r *http.Request) error {
	payload, err := json.Marshal(
		GetModel{
			ID: "123456",
		},
	)
	if err != nil {
		return err
	}

	if err := s.Event.Producer(
		GATEWAY_EVENT_GET,
		payload,
	); err != nil {
		return err
	}
	return WriteJson(w, http.StatusOK, "in process get")
}

type createHandlerFunc func(w http.ResponseWriter, r *http.Request) error

func handleFunc(f createHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJson(w, http.StatusBadRequest, fmt.Errorf("error in gateway: %s", err))
			return
		}
	}
}

func WriteJson(w http.ResponseWriter, code int, msg any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(msg)
}
