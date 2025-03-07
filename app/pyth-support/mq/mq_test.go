package mq_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"

	"pyth-go/app/pyth-support/mq/consumer"
	"pyth-go/app/pyth-support/mq/producer"
)

var groupCfg = consumer.ConsumerGroupConfig{
	Addrs:          []string{"localhost:9194", "localhost:9294", "localhost:9394"},
	GroupID:        "test-group",
	GroupTopics:    []string{"test-topic"},
	ConsumerNumber: 3,
}

func TestStartConsumer(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	// init producer and consumer
	producer := producer.NewProducer(ctx, groupCfg.Addrs, "test-topic")
	consumer := consumer.NewConsumer(ctx, groupCfg)

	// expect x messages
	msgNum := 10
	var wg sync.WaitGroup
	wg.Add(msgNum)

	// store messages
	var mu sync.Mutex
	var consumedMessages []string

	// hanlder function
	fn := func(m kafka.Message) error {
		mu.Lock()
		consumedMessages = append(consumedMessages, string(m.Value))
		mu.Unlock()
		wg.Done()
		return nil
	}

	// start consumer
	go consumer.Start(ctx, fn)

	// produce x messages
	msgs := make([]kafka.Message, msgNum)
	for i := range msgNum {
		msg := kafka.Message{
			Value: []byte(fmt.Sprintf("Test message %d", i)),
		}
		msgs[i] = msg
	}
	if err := producer.Do(ctx, msgs); err != nil {
		t.Fatalf("failed to produce message: %v", err)
	}

	done := make(chan struct{})

	go func() {
		wg.Wait()
		time.Sleep(1 * time.Second)
		close(done)
	}()

	select {
	case <-done:
		// all messages consumed
	case <-ctx.Done():
		t.Error("test timed out waiting for messages")
	}

	// check whether # of messages received is correct
	if len(consumedMessages) != msgNum {
		t.Errorf("expected 10 messages, got %d", len(consumedMessages))
	}

	// Check whether messages ordering
	expected := make(map[string]bool)
	for i := 0; i < msgNum; i++ {
		expected[fmt.Sprintf("Test message %d", i)] = true
	}
	for _, m := range consumedMessages {
		if !expected[m] {
			t.Errorf("unexpected message: %s", m)
		}
	}
}

func TestStartConsumerGroup(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	// init producer and consumer
	producer := producer.NewProducer(ctx, groupCfg.Addrs, "test-topic")
	consumerGroup := consumer.NewConsumerGroup(ctx, groupCfg)

	msgNum := 10
	// expect x messages
	var wg sync.WaitGroup
	wg.Add(msgNum)

	// store messages
	var mu sync.Mutex
	var consumedMessages []string

	// hanlder function
	fn := func(m kafka.Message) error {
		mu.Lock()
		consumedMessages = append(consumedMessages, string(m.Value))
		mu.Unlock()
		wg.Done()
		return nil
	}

	// start consumer
	go consumerGroup.Start(ctx, fn)

	// produce x messages
	msgs := make([]kafka.Message, msgNum)
	for i := range msgNum {
		msg := kafka.Message{
			Value: []byte(fmt.Sprintf("Test message %d", i)),
		}
		msgs[i] = msg
	}
	if err := producer.Do(ctx, msgs); err != nil {
		t.Fatalf("failed to produce message: %v", err)
	}

	done := make(chan struct{})

	go func() {
		wg.Wait()
		time.Sleep(1 * time.Second)
		close(done)
	}()

	select {
	case <-done:
		// all messages consumed
	case <-ctx.Done():
		t.Error("test timed out waiting for messages")
	}

	// check whether # of messages received is correct
	if len(consumedMessages) != msgNum {
		t.Errorf("expected 10 messages, got %d", len(consumedMessages))
	}

	// Check whether messages ordering
	expected := make(map[string]bool)
	for i := 0; i < msgNum; i++ {
		expected[fmt.Sprintf("Test message %d", i)] = true
	}
	for _, m := range consumedMessages {
		if !expected[m] {
			t.Errorf("unexpected message: %s", m)
		}
	}
}
