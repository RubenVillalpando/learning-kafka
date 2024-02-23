package main

import (
	"github.com/RubenVillalpando/learning-kafka/internal/application"
	"github.com/RubenVillalpando/learning-kafka/internal/kafka"
)

func main() {
	app := application.New()

	app.ServeHTTP()
	// franz go
	kafka.ProduceMessage()
}
