package main

import (
	"github.com/RubenVillalpando/learning-kafka/internal/application"
)

func main() {
	app := application.New()

	app.StartKafkaConsumer()
	app.Serve()
}
