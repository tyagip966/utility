package kafka

type KafkaProducerConfig struct {
	Brokers      []string `validate:"gt=0,dive,required"`
	ClientId     string   `validate:"required"`
	KafkaVersion string   `validate:"required,validKafkaVersion"`
}
