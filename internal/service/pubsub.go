package service

import (
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"os"
	"sweng-task/internal/config"
)

type PubSub struct {
	logs       *zap.SugaredLogger
	cfg        *config.Config
	connection sarama.SyncProducer
}

func NewPubSub(log *zap.SugaredLogger, cfg *config.Config) *PubSub {
	return &PubSub{
		logs: log,
		cfg:  cfg,
	}
}

func (p *PubSub) Connect(config *sarama.Config) error {
	brokers := []string{os.Getenv("BROKER")}
	producer, err := sarama.NewSyncProducer([]string{p.cfg.PubSub.Broker}, config)
	if err != nil {
		p.logs.Fatalw("Failed to create Kafka producer",
			"error", err,
			"brokers", brokers,
		)
	} else {
		p.logs.Info("Kafka producer connected to %v", brokers)
	}
	p.connection = producer
	return nil
}

func (p *PubSub) Connection() sarama.SyncProducer {
	return p.connection
}

func (p *PubSub) Publish(data string) {
	_, _, err := p.connection.SendMessage(&sarama.ProducerMessage{
		Topic: p.cfg.PubSub.Topic,
		Value: sarama.StringEncoder(data),
	})
	if err != nil {
		p.logs.Warnf(err.Error())
	}
}
