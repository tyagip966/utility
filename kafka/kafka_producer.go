package kafka

import (
	"github.com/Shopify/sarama"
)

type KafkaProducer interface {
	SendMessage(topic string, message interface{}) (partition int32, offset int64, err error)
	SendKeyAndMessage(topic string, messageKey string, message interface{}) (partition int32, offset int64, err error)
	Close() error
}

type kafkaProducer struct {
	syncProducer sarama.SyncProducer
}

func (producer kafkaProducer) Close() error {
	err := producer.syncProducer.Close()
	return err
}

func NewKafkaProducer(producer sarama.SyncProducer) KafkaProducer {
	kafkaProducer := kafkaProducer{syncProducer: producer}
	return &kafkaProducer
}

func (producer kafkaProducer) SendMessage(topic string, message interface{}) (partition int32, offset int64, err error) {
	encoder := NewJsonEncoder(message)
	producerMessage := sarama.ProducerMessage{Topic: topic, Value: encoder}

	partition, offset, err = producer.syncProducer.SendMessage(&producerMessage)

	return partition, offset, err
}

func (producer kafkaProducer) SendKeyAndMessage(topic string, messageKey string, message interface{}) (partition int32, offset int64, err error) {
	encoder := NewJsonEncoder(message)
	producerMessage := sarama.ProducerMessage{Topic: topic, Value: encoder, Key: sarama.StringEncoder(messageKey)}

	partition, offset, err = producer.syncProducer.SendMessage(&producerMessage)

	return partition, offset, err
}
