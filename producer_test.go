package main

import (
	"context"
	"testing"
	"time"
)

var topic = "epoch"

func TestSendEpochMessage(t *testing.T) {

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	producer, err := newProducer()
	if err != nil {
		t.Fatal()
	}
	err = sendEpochMessage(ctx, producer, topic)
	if err != nil {
		t.Error()
	}

}
