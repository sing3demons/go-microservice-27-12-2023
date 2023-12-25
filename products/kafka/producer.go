package kafka

import (
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

func NewSyncProducer(kafkaBrokers []string) (EventProducer, error) {
	producer, err := sarama.NewSyncProducer(kafkaBrokers, nil)
	if err != nil {
		return nil, err
	}

	return &eventProducer{producer}, nil
}

type EventProducer interface {
	Produce(topic string, event any) (err error)
	Close() error
}

type eventProducer struct {
	producer sarama.SyncProducer
}

func (e *eventProducer) Close() error {
	return e.producer.Close()
}

func (e *eventProducer) Produce(topic string, event any) (err error) {
	value, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(value),
	}
	partition, offset, err := e.producer.SendMessage(&msg)
	if err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"topic":     topic,
		"partition": partition,
		"offset":    offset,
		"event":     event,
	}).Info("PRODUCE_EVENT_SUCCESS")

	return nil
}
