package main

import (
	"context"
	"os"
	"os/signal"
	exporter "producer/metrics"
	"strconv"
	"time"

	cmap "github.com/orcaman/concurrent-map"
	log "github.com/sirupsen/logrus"

	"github.com/Shopify/sarama"
)

func init() {
	log.SetLevel(logLevel())
}

func main() {
	topic := "epoch"
	if value, ok := os.LookupEnv("TOPIC"); ok {
		topic = value
	}
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	go exporter.Exporter()
	go func() {
		oscall := <-c
		log.Debug("system call:%+v", oscall)
		cancel()
	}()

	producer, err := newProducer()
	if err != nil {
		log.Error("Could not create producer: ", err)
	}

	if err := sendMessage(ctx, producer, topic); err != nil {
		log.Error("failed to produce:+%v\n", err)
	}
}

func sendMessage(ctx context.Context, producer sarama.SyncProducer, topic string) (err error) {
	log.Info("Server Start Producing")
	var partitionProduced = cmap.New()
	go func() {
		for {
			time.Sleep(time.Millisecond)
			epochTime := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
			msg := prepareMessage(topic, epochTime)
			partition, offset, err := producer.SendMessage(msg)
			if err != nil {
				log.Panic("error occured.", err)
			}
			counter, _ := partitionProduced.Get(string(partition))
			if counter == nil {
				partitionProduced.Set(string(partition), 1)
			} else {
				counted := counter.(int)
				counted++
				partitionProduced.Set(string(partition), counted)
				log.Trace("Counter: ", " Message: ", epochTime, " topic: ", topic, " partition: ", partition, " offset: ", offset)
				M := exporter.ProducedMessageCounter.WithLabelValues(strconv.Itoa(int(partition)), topic)
				M.Inc()
			}
		}
	}()
	<-ctx.Done()
	log.Info("Server Stopped")

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	log.Info("Server Exited Properly")

	return
}

func prepareMessage(topic, message string) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: -1,
		Value:     sarama.StringEncoder(message),
	}

	return msg
}
