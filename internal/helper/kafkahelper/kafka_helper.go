package kafkahelper

import (
	"time"

	"github.com/IBM/sarama"
)

type KafkaProducer struct {
	Producer sarama.SyncProducer
	Broker   []string
}

func NewKafkaProducer(brokers []string) (*KafkaProducer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 3
	cfg.Producer.Retry.Backoff = 100 * time.Millisecond
	cfg.Version = sarama.V3_0_2_0

	producer, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{
		Producer: producer,
		Broker:   brokers,
	}, nil
}

func (kp *KafkaProducer) Publish(topic, key string, value []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(value),
	}

	_, _, err := kp.Producer.SendMessage(msg)
	return err
}

func (kp *KafkaProducer) Close() error {
	return kp.Producer.Close()
}
