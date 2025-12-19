package kafkahelper

import (
	"log"
	"os"
	"time"

	"github.com/IBM/sarama"
)

type KafkaProducer struct {
	Producer sarama.SyncProducer
	Brokers  []string
}

func NewKafkaProducer(brokers []string) (*KafkaProducer, error) {
	cfg := sarama.NewConfig()

	// producer config
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 3
	cfg.Producer.Retry.Backoff = 100 * time.Millisecond
	cfg.Producer.Timeout = 5 * time.Second
	cfg.Producer.Partitioner = sarama.NewHashPartitioner

	// network config
	cfg.Net.DialTimeout = 5 * time.Second
	cfg.Net.ReadTimeout = 5 * time.Second
	cfg.Net.WriteTimeout = 5 * time.Second

	cfg.Version = sarama.V3_0_2_0

	sarama.Logger = log.New(os.Stdout, "[Sarama]", log.LstdFlags)

	producer, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{
		Producer: producer,
		Brokers:  brokers,
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
	if kp.Producer == nil {
		return nil
	}
	return kp.Producer.Close()
}
