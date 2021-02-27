package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/BurntSushi/toml"
	"github.com/Shopify/sarama"
)

type config struct {
	KafkaWorkers   int       `toml:"kafka_output_workers"`
	KafkaBrokers   []string  `toml:"kafka_brokers"`
	KafkaTopic     string    `toml:"kafka_topic"`
	KafkaBatchSize int       `toml:"kafka_batch_size"`
	KafkaQueueSize int       `toml:"kafka_queue_size"`
	KafkaFlushTime int       `toml:"kafka_flush_time_ms"`
	LogLevel       log.Level `toml:"loglevel"`
}

func main() {
	var conf config
	_, err := toml.DecodeFile("./conf.toml", &conf)

	if err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		oscall := <-c
		log.Printf("system call:%+v", oscall)
		cancel()
	}()

	producer, err := newProducer(&conf)
	if err != nil {
		fmt.Println("Could not create producer: ", err)
	}

	if err := sendMessage(ctx, producer, conf.KafkaTopic); err != nil {
		log.Printf("failed to produce:+%v\n", err)
	}
}

func sendMessage(ctx context.Context, producer sarama.SyncProducer, topic string) (err error) {
	//var wg sync.WaitGroup
	go func() {
		for {
			time.Sleep(time.Millisecond)
			msg := prepareMessage(topic, strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10))
			_, _, err := producer.SendMessage(msg)
			if err != nil {
				fmt.Println("%s error occured.", err.Error())
			}
		}
	}()
	<-ctx.Done()

	log.Printf("server stopped")

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	log.Printf("server exited properly")

	return
}

func newProducer(conf *config) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Flush.Frequency = time.Duration(conf.KafkaFlushTime) * time.Millisecond
	config.Producer.Flush.Messages = conf.KafkaBatchSize
	producer, err := sarama.NewSyncProducer(conf.KafkaBrokers, config)

	return producer, err
}

func prepareMessage(topic, message string) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: -1,
		Value:     sarama.StringEncoder(message),
	}

	return msg
}