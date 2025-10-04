package mq

import (
	"context"
	"fmt"
	"time"

	"github.com/NUS-ISS-Agile-Team/ceramicraft-user-mservice/server/config"
	"github.com/NUS-ISS-Agile-Team/ceramicraft-user-mservice/server/log"
	"github.com/segmentio/kafka-go"
)

var (
	kafkaProducerImpl KafkaProducer
)

type KafkaProducer interface {
	Produce(ctx context.Context, topic string, key string, value []byte) error
}

type KafkaProducerImpl struct {
	producer *kafka.Writer
}

func InitKafka() {
	producer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      config.Config.KafkaConfig.Brokers,
		Balancer:     &kafka.Hash{},
		BatchSize:    config.Config.KafkaConfig.BatchSize,
		BatchBytes:   config.Config.KafkaConfig.MaxBytes,
		BatchTimeout: time.Duration(config.Config.KafkaConfig.BatchTimeoutMillis) * time.Millisecond,
		RequiredAcks: config.Config.KafkaConfig.Acks,
		MaxAttempts:  config.Config.KafkaConfig.Retries,
	})

	kafkaProducerImpl = &KafkaProducerImpl{producer: producer}
	log.Logger.Infof("Kafka producer initialized")
}

func GetKafkaProducer() KafkaProducer {
	return kafkaProducerImpl
}

func (k *KafkaProducerImpl) Produce(ctx context.Context, topic string, key string, value []byte) error {
	if value == nil {
		return fmt.Errorf("value cannot be nil")
	}

	err := k.producer.WriteMessages(ctx, kafka.Message{Topic: topic, Key: []byte(key), Value: value})
	log.Logger.Infof("Produced message to topic %s: key=%s, value=%s, err=%v", topic, key, string(value), err)
	return err
}
