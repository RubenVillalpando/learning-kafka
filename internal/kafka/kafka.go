package kafka

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RubenVillalpando/learning-kafka/internal/model"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Client struct {
	producer *kgo.Client
	consumer *kgo.Client
}

func New() *Client {
	p, err := kgo.NewClient(
		kgo.WithLogger(kgo.BasicLogger(os.Stderr, kgo.LogLevelInfo, func() string {
			return time.Now().Format("[2006-01-02 15:04:05.999] ")
		})),
		kgo.DefaultProduceTopic("customer-info"),
		kgo.SeedBrokers("localhost:9092"),
	)
	if err != nil {
		panic(err)
	}
	c, err := kgo.NewClient(
		kgo.WithLogger(kgo.BasicLogger(os.Stderr, kgo.LogLevelInfo, func() string {
			return time.Now().Format("[2006-01-02 15:04:05.999] ")
		})),
		kgo.ConsumeTopics("customer-info"),
		kgo.SeedBrokers("localhost:9092"),
	)
	if err != nil {
		panic(err)
	}
	return &Client{
		producer: p,
		consumer: c,
	}
}

func (c *Client) ProduceMessage(customer *model.Customer) error {

	ctx := context.Background()

	record := &kgo.Record{Topic: "customer-info", Value: customer.Value(), Key: customer.Key()}

	c.producer.Produce(ctx, record, func(r *kgo.Record, err error) {
		if err != nil {
			log.Printf("Record with key {%v} had a produce error: %v\n", r.Key, err)
		}
		log.Printf("key %s value %s\n", r.Key, r.Value)
	})

	return nil
}

func (c *Client) Poll() {
	log.Println("before starting gofunc")

	ctx := context.Background()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		defer func() {
			c.consumer.Close()
			if r := recover(); r != nil {
				log.Println("Consumer stopped: ", r)
				time.Sleep(5 * time.Second)
				c.Poll()
			}
		}()
		for {
			select {
			case <-sigChan:
				fmt.Println("Received Signal from terminal, shutting down consumer...")
				return
			default:
				fetches := c.consumer.PollFetches(ctx)
				if errs := fetches.Errors(); len(errs) > 0 {
					panic(fmt.Sprint(errs))
				}
				iter := fetches.RecordIter()
				for !iter.Done() {
					record := iter.Next()
					log.Println(string(record.Value), "from an iterator!")
				}
			}
		}
	}()
}
