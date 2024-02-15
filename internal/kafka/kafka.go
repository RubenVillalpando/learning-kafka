package kafka

import (
	"context"
	"fmt"
	"sync"

	"github.com/RubenVillalpando/learning-kafka/internal/model"
	"github.com/twmb/franz-go/pkg/kgo"
)

func ProduceMessage() {
	seeds := []string{"localhost:9092"}

	customer := &model.Customer{
		ID:    0,
		Name:  "Ruben Villalpando",
		Email: "rubenvillalpando1299@gmail.com",
	}
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
	)
	println("1")
	if err != nil {
		panic(err)
	}
	println("2")
	defer cl.Close()

	ctx := context.Background()

	val, err := customer.Value()
	if err != nil {
		fmt.Println("error marshaling val: ", err)
		return
	}
	println("3")

	var wg sync.WaitGroup
	wg.Add(1)
	record := &kgo.Record{Topic: "customer-info", Value: val, Key: customer.Key()}
	println("4")

	cl.Produce(ctx, record, func(r *kgo.Record, err error) {
		println("5")
		defer wg.Done()
		if err != nil {
			fmt.Printf("Record with key {%v} had a produce error: %v\n", r.Key, err)
		}
		fmt.Printf("key %v value %v", r.Key, r.Value)
	})
	println("7")
	wg.Wait()
	println("6")

}
