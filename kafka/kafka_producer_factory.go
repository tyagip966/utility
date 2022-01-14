package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/go-playground/validator/v10"
)

type KafkaProducerFactory interface {
	NewKafkaProducer(config KafkaProducerConfig) (KafkaProducer, error)
}

type kafkaProducerFactory struct {
}

func NewKafkaProducerFactory() KafkaProducerFactory {
	return &kafkaProducerFactory{}
}

func (factory kafkaProducerFactory) NewKafkaProducer(producerConfig KafkaProducerConfig) (KafkaProducer, error) {

	syncProducerConfig, configError := factory.buildSyncProducerConfigFrom(producerConfig)
	if configError != nil {
		return nil, configError
	}

	syncProducer, producerError := sarama.NewSyncProducer(producerConfig.Brokers, syncProducerConfig)
	if producerError != nil {
		return nil, producerError
	}

	return NewKafkaProducer(syncProducer), nil
}

func (factory kafkaProducerFactory) buildSyncProducerConfigFrom(config KafkaProducerConfig) (*sarama.Config, error) {
	if validationError := factory.validateConfig(config); validationError != nil {
		return nil, validationError
	}

	var kafkaVersion sarama.KafkaVersion
	var kafkaVersionError error
	if kafkaVersion, kafkaVersionError = sarama.ParseKafkaVersion(config.KafkaVersion); kafkaVersionError != nil {
		return nil, kafkaVersionError
	}

	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.ClientID = config.ClientId
	saramaConfig.Version = kafkaVersion

	return saramaConfig, nil
}

func (factory kafkaProducerFactory) validateConfig(config KafkaProducerConfig) error {
	validate := validator.New()
	if validationError := validate.RegisterValidation("validKafkaVersion", IsValidKafkaVersion); validationError != nil {
		return validationError
	}
	return validate.Struct(config)
}
