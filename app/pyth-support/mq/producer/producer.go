package producer

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"

	
)

type MQProducer struct {
	Writer *kafka.Writer
}

func NewProducer(ctx context.Context, addrs []string, topic string) MQProducer {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: addrs,
		Topic:   topic,
	})
	log.Info().Msgf("Producer created for Topic {%s}", topic)
	return MQProducer{
		Writer: w,
	}
}

func (p *MQProducer) Do(ctx context.Context, msgs []kafka.Message) error {
	return p.Writer.WriteMessages(ctx, msgs...)
}
