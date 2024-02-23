package kafka

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/RubenVillalpando/learning-kafka/internal/model"
	"github.com/twmb/franz-go/pkg/kgo"
)

type KafkaClient struct {
	Client *kgo.Client
}

func New() *KafkaClient {
	cl, err := kgo.NewClient(
		kgo.WithLogger(kgo.BasicLogger(os.Stderr, kgo.LogLevelInfo, func() string {
			return time.Now().Format("[2006-01-02 15:04:05.999] ")
		})),
		kgo.ConsumeTopics("consumer-info"),
		kgo.DefaultProduceTopic("consumer-info"),
		kgo.SeedBrokers("localhost:9092"),
	)
	if err != nil {
		panic(err)
	}
	return &KafkaClient{
		Client: cl,
	}
}

func (kc *KafkaClient) ProduceMessage(m *model) {

	customer := &model.Customer{
		ID:    0,
		Name:  "Ruben Villalpando",
		Email: "rubenvillalpando1299@gmail.com",
	}

	ctx := context.Background()

	val, err := customer.Value()
	if err != nil {
		fmt.Println("error marshaling val: ", err)
		return
	}

	var wg sync.WaitGroup

	for i := 0; i < 1; i++ {
		wg.Add(1)
		record := &kgo.Record{Topic: "customer-info", Value: val, Key: customer.Key()}

		cl.Produce(ctx, record, func(r *kgo.Record, err error) {
			defer wg.Done()
			if err != nil {
				fmt.Printf("Record with key {%v} had a produce error: %v\n", r.Key, err)
			}
			fmt.Printf("key %v value %v\n", r.Key, r.Value)
		})
	}

	wg.Wait()
	for {
		fetches := cl.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			panic(fmt.Sprint(errs))
		}

		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()
			fmt.Println(string(record.Value), "from an iterator!")
		}
	}

}
