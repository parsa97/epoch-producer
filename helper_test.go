package main

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestNewProducer(t *testing.T) {
	os.Setenv("PRODUCER_MAX_MESSAGE_BYTES", "1000")
	os.Setenv("PRODUCER_FLUSH_FREQUENCY", "10")
	os.Setenv("PRODUCER_FLUSH_MESSAGE", "10")
	os.Setenv("PRODUCER_FLUSH_MAX_MESSAGE", "10")
	os.Setenv("PRODUCER_RETURN_SUCCESS", "true")
	os.Setenv("PRODUCER_TIMEOUT", "50")
	os.Setenv("PRODUCER_RETRY_MAX", "3")
	os.Setenv("PRODUCER_RETRY_BACKOFF", "5")
	os.Setenv("PRODUCER_RETURN_ERROR", "true")
	os.Setenv("PRODUCER_COMPRESSIONLEVEL", "0")
	os.Setenv("PRODUCER_CLIENTID", "sarama_test")
	os.Setenv("PRODUCER_CHANNELBUFFERSIZE", "50")
	os.Setenv("PRODUCER_VERSION", "2.6.0")
	os.Setenv("PRODUCER_BROKERS", "localhost:9092")

	_, err := newProducer()
	if err != nil {
		t.Error()
	}
}

func TestRequiredAcks(t *testing.T) {
	os.Setenv("PRODUCER_REQUIRED_ACKS", "-1")
	response := saramaPartitioner()
	if response == nil {
		t.Error()
	}
}

func TestSaramaPartitioner(t *testing.T) {
	os.Setenv("PRODUCER_PARTITIONER", "badvalue")
	response := saramaPartitioner()
	if response == nil {
		t.Error()
	}
}

func TestLogLevel(t *testing.T) {
	os.Setenv("LOG_LEVEL", "badvalue")
	response := logLevel()
	if response != log.InfoLevel {
		t.Error()
	}
}
