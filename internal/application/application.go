package application

import (
	"log"
	"net/http"

	"github.com/RubenVillalpando/learning-kafka/internal/handler"
	"github.com/RubenVillalpando/learning-kafka/internal/kafka"
)

type App struct {
	consumerHandler *handler.Handler
	kafkaClient     *kafka.Client
}

func New() *App {
	kClient := kafka.New()
	cHandler := handler.New(kClient)
	return &App{
		consumerHandler: cHandler,
		kafkaClient:     kClient,
	}
}

func (app *App) Serve() error {
	http.Handle("/customer", app.consumerHandler)

	// Start the HTTP server on port 8080
	log.Println("Starting server on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
		return err
	}
	return nil
}

func (app *App) StartKafkaConsumer() {
	app.kafkaClient.Poll()
}
