package main

import (
	"context"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

var topic = "test"

func TestSendEpochMessage(t *testing.T) {

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	producer, err := newProducer()
	if err != nil {
		log.Error("Could not create producer: ", err)
	}
	err = sendEpochMessage(ctx, producer, topic)
	if err == nil {
		log.Info("Everything works well")
	}

}
