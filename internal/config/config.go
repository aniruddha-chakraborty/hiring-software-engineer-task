package config

import (
	"github.com/IBM/sarama"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config represents the application configuration
type Config struct {
	App     AppConfig     `split_words:"true"`
	Server  ServerConfig  `split_words:"true"`
	PubSub  PubSubConfig  `split_words:"true"`
	Metrics MetricsConfig `split_words:"true"`
}

// AppConfig contains application-specific configuration
type AppConfig struct {
	Name        string `default:"Ad Bidding Service"`
	Environment string `default:"development"`
	LogLevel    string `default:"info" split_words:"true"`
	Version     string `default:"1.0.0"`
}

// ServerConfig contains HTTP server configuration
type ServerConfig struct {
	Port    int           `default:"8080"`
	Timeout time.Duration `default:"30s"`
}

// PubSubConfig contains connection requirement
type PubSubConfig struct {
	Broker        string `default:"kafka:9092"`
	Topic         string `default:"tracking-events"`
	RetryMax      int    `default:"3"`
	ReturnSuccess bool   `default:"true"`
}

type MetricsConfig struct {
	Port int `default:"9100"`
}

//Kafka config spin up

func KafkaConfigLoad() *sarama.Config {
	p := PubSubConfig{}
	KafkaConfig := sarama.NewConfig()
	KafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	KafkaConfig.Producer.Retry.Max = p.RetryMax
	KafkaConfig.Producer.Return.Successes = true
	return KafkaConfig
}

// Load loads the configuration from environment variables
func Load() (*Config, error) {
	var config Config
	if err := envconfig.Process("app", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
