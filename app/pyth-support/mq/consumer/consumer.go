package consumer

import (
	"context"

	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
)

type MQConsumer struct {
	Reader *kafka.Reader
	ID     int
}

type ConsumerGroupConfig struct {
	Addrs          []string
	GroupTopics    []string
	GroupID        string
	ConsumerNumber int
}

func NewConsumer(groupCfg ConsumerGroupConfig) MQConsumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     groupCfg.Addrs,
		GroupTopics: groupCfg.GroupTopics,
		GroupID:     groupCfg.GroupID,
		MinBytes:    10e3,
		MaxBytes:    10e3,
	})
	logx.Infof("Consumer created for topic {%+v}", groupCfg.GroupTopics)
	return MQConsumer{
		Reader: r,
	}
}

func (c *MQConsumer) Start(ctx context.Context, fn func(m kafka.Message) error) {
	r := c.Reader
	for {
		select {
		case <-ctx.Done():
			logx.Info("MQConsumer: context cancelled, exiting consumer loop")
			return
		default:
		}
		m, err := r.FetchMessage(ctx)
		if err != nil {
			logx.Errorw("MQConsumer: failed to fetch message ", logx.Field("error", err))
			continue
		}
		if err := fn(m); err != nil {
			logx.Errorw("MQConsumer: failed to process message ", logx.Field("error", err))
			continue
		}
		if err := r.CommitMessages(ctx, m); err != nil {
			logx.Errorw("MQConsumer: failed to commit message ", logx.Field("error", err))
			continue
		}
		logx.Infow("MQConsumer: message consumed ",
			logx.Field("topic", m.Topic),
			logx.Field("id", c.ID),
			logx.Field("partition", m.Partition),
			logx.Field("offset", m.Offset),
			logx.Field("key", string(m.Key)),
			logx.Field("value", string(m.Value)),
		)
	}
}

func (c *MQConsumer) Stop() error {
	return c.Reader.Close()
}

type MQConsumerGroup struct {
	Consumers []*MQConsumer
	GroupCfg  ConsumerGroupConfig
}

func NewConsumerGroup(groupCfg ConsumerGroupConfig) MQConsumerGroup {
	consumers := make([]*MQConsumer, groupCfg.ConsumerNumber)
	for i := 0; i < groupCfg.ConsumerNumber; i++ {
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers:     groupCfg.Addrs,
			GroupTopics: groupCfg.GroupTopics,
			GroupID:     groupCfg.GroupID,
		})
		consumer := &MQConsumer{
			Reader: r,
			ID:     i,
		}
		consumers[i] = consumer
	}
	logx.Infow("ConsumerGroup created new topic ", logx.Field("topics", groupCfg.GroupTopics))
	return MQConsumerGroup{
		Consumers: consumers,
		GroupCfg:  groupCfg,
	}
}

func (cg *MQConsumerGroup) Start(ctx context.Context, fn func(m kafka.Message) error) {
	for _, c := range cg.Consumers {
		go c.Start(ctx, fn)
	}
}

func (cg *MQConsumerGroup) Stop() {
	for _, c := range cg.Consumers {
		err := c.Stop()
		if err != nil {
			logx.Errorw("MQCOnsumerGroup close consumer ", logx.Field("error", err))
		}
	}
}
