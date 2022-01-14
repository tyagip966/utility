package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/go-playground/validator/v10"
)

func IsValidKafkaVersion(fieldLevel validator.FieldLevel) bool {
	kafkaVersion := fieldLevel.Field().String()
	_, err := sarama.ParseKafkaVersion(kafkaVersion)

	if err != nil {
		return false
	}
	return true
}
