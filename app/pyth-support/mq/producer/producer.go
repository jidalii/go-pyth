package producer

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

type MQProducer struct {
	Writer *kafka.Writer
}

type MQProducerConf struct {
	Addrs []string
	Topic string
}

func NewProducer(cfg MQProducerConf) MQProducer {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: cfg.Addrs,
		Topic:   cfg.Topic,
	})
	log.Info().Msgf("Producer created for Topic {%s}", cfg.Topic)
	return MQProducer{
		Writer: w,
	}
}

func (p *MQProducer) Start() {
}

func (p *MQProducer) Do(ctx context.Context, msgs ...kafka.Message) error {
	return p.Writer.WriteMessages(ctx, msgs...)
}

func (p *MQProducer) Stop() error {
	return p.Writer.Close()
}
