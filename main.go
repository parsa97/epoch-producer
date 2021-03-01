package main

import (
	"context"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"time"

	//app "consumer/app"

	log "github.com/sirupsen/logrus"

	"github.com/Shopify/sarama"
)

func logLevel() log.Level {
	if value, ok := os.LookupEnv("LOG_LEVEL"); ok {
		switch value {
		case "trace":
			return log.TraceLevel
		case "debug":
			return log.DebugLevel
		case "info":
			return log.InfoLevel
		case "warn":
			return log.WarnLevel
		case "error":
			return log.ErrorLevel
		case "fatal":
			return log.FatalLevel
		case "panic":
			return log.WarnLevel
		default:
			return log.InfoLevel
		}
	}
	return log.InfoLevel
}

func main() {
	log.SetLevel(logLevel())
	topic := "epoch"
	if value, ok := os.LookupEnv("TOPIC"); ok {
		topic = value
	}
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

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
	var counter int64

	go func() {
		for {
			time.Sleep(time.Millisecond)
			rawString := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
			msg := prepareMessage(topic, rawString)
			partition, offset, err := producer.SendMessage(msg)
			if err != nil {
				log.Panic("error occured.", err)
			}
			log.Trace("Counter: ", counter, " Message: ", rawString, " topic: ", topic, " partition: ", partition, " offset: ", offset)
			counter++
		}
	}()
	<-ctx.Done()
	log.Info("messages produced: ", counter)
	log.Info("server stopped")

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	log.Info("server exited properly")

	return
}
func requiredAcks() sarama.RequiredAcks {
	if value, ok := os.LookupEnv("PRODUCER_REQUIRED_ACKS"); ok {
		valuei, err := strconv.Atoi(value)
		if err != nil {
			log.Error("Bad! PRODUCER_REQUIRED_ACKS: ", err)
		}
		switch valuei {
		case 0:
			return sarama.NoResponse
		case 1:
			return sarama.WaitForLocal
		case -1:
			return sarama.WaitForAll
		default:
			return sarama.WaitForLocal
		}
	}
	return sarama.WaitForLocal
}

func saramaPartitioner() sarama.PartitionerConstructor {
	if value, ok := os.LookupEnv("PRODUCER_PARTITIONER"); ok {
		switch value {
		case "random":
			log.Debug("Partitioner: Random")
			return sarama.NewRandomPartitioner
		case "hash":
			log.Debug("Partitioner: Hash")
			return sarama.NewHashPartitioner
		case "rr":
			log.Debug("Partitioner: RoundRobin")
			return sarama.NewRoundRobinPartitioner
		default:
			log.Debug("Partitioner: Hash")
			return sarama.NewHashPartitioner
		}
	}
	return sarama.NewHashPartitioner
}
func newProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	var brokers []string
	brokers = []string{"localhost:9092"}

	config.Producer.RequiredAcks = requiredAcks()
	if value, ok := os.LookupEnv("PRODUCER_MAX_MESSAGE_BYTES"); ok {
		valuei, err := strconv.Atoi(value)
		if err != nil {
			log.Error("Bad! PRODUCER_MAX_MESSAGE_BYTES: ", err)
		}
		config.Producer.MaxMessageBytes = valuei
	}
	if value, ok := os.LookupEnv("PRODUCER_FLUSH_FREQUENCY"); ok {
		valuei, err := strconv.Atoi(value)
		if err != nil {
			log.Error("Bad! PRODUCER_FLUSH_FREQUENCY: ", err)
		}
		config.Producer.Flush.Frequency = time.Duration(rand.Int31n(int32(valuei))) * time.Millisecond
	}
	if value, ok := os.LookupEnv("PRODUCER_FLUSH_MESSAGE"); ok {
		valuei, err := strconv.Atoi(value)
		if err != nil {
			log.Error("Bad! PRODUCER_FLUSH_MESSAGE: ", err)
		}
		config.Producer.Flush.Messages = valuei
	}
	if value, ok := os.LookupEnv("PRODUCER_FLUSH_MAX_MESSAGE"); ok {
		valuei, err := strconv.Atoi(value)
		if err != nil {
			log.Error("Bad! PRODUCER_FLUSH_MAX_MESSAGE: ", err)
		}
		config.Producer.Flush.MaxMessages = valuei
	}
	config.Producer.Return.Successes = true
	if value, ok := os.LookupEnv("PRODUCER_RETURN_SUCCESS"); ok {
		if value == "false" {
			config.Producer.Return.Successes = false
		}
	}
	if value, ok := os.LookupEnv("PRODUCER_TIMEOUT"); ok {
		valuei, err := strconv.Atoi(value)
		if err != nil {
			log.Error("Bad! PRODUCER_TIMEOUT: ", err)
		}
		config.Producer.Timeout = time.Duration(rand.Int31n(int32(valuei))) * time.Second
	}
	config.Producer.Partitioner = saramaPartitioner()

	if value, ok := os.LookupEnv("PRODUCER_RETRY_MAX"); ok {
		valuei, err := strconv.Atoi(value)
		if err != nil {
			log.Error("Bad! PRODUCER_RETRY_MAX: ", err)
		}
		config.Producer.Retry.Max = int(valuei)
	}
	if value, ok := os.LookupEnv("PRODUCER_RETRY_BACKOFF"); ok {
		valuei, err := strconv.Atoi(value)
		if err != nil {
			log.Error("Bad! PRODUCER_RETRY_BACKOFF: ", err)
		}
		config.Producer.Retry.Backoff = time.Duration(rand.Int31n(int32(valuei))) * time.Millisecond
	}
	if value, ok := os.LookupEnv("PRODUCER_RETURN_ERROR"); ok {
		if value == "true" {
			config.Producer.Return.Errors = true
		} else {
			config.Producer.Return.Errors = false
		}
	}
	if value, ok := os.LookupEnv("PRODUCER_COMPRESSIONLEVEL"); ok {
		valuei, err := strconv.Atoi(value)
		if err != nil {
			log.Error("Bad! PRODUCER_COMPRESSIONLEVEL: ", err)
		}
		config.Producer.CompressionLevel = valuei
	}
	if value, ok := os.LookupEnv("PRODUCER_CLIENTID"); ok {
		config.ClientID = value
	}
	if value, ok := os.LookupEnv("PRODUCER_CHANNELBUFFERSIZE"); ok {
		valuei, err := strconv.Atoi(value)
		if err != nil {
			log.Error("Bad! PRODUCER_CHANNELBUFFERSIZE: ", err)
		}
		config.ChannelBufferSize = valuei
	}
	if value, ok := os.LookupEnv("PRODUCER_VERSION"); ok {
		version, err := sarama.ParseKafkaVersion(value)
		if err != nil {
			log.Error("version incompatible: ", err)
		}
		config.Version = version
	}
	if value, ok := os.LookupEnv("PRODUCER_BROKERS"); ok {
		brokers = []string{value}
	}
	producer, err := sarama.NewSyncProducer(brokers, config)
	log.Info(config.Producer.Partitioner)
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
