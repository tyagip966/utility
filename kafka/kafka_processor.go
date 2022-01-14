package kafka

import (
	"context"
	"errors"
	"github.com/Shopify/sarama"
	"log"
	"strings"
	"time"
)

type KafkaUtil interface {
	StartKafka(retryDelay int)
	GetKafkaOffSet() (int64, error)
}

const (
	OffsetOldest = "oldest"
	OffsetNewest = "newest"
)

type kafkaUtil struct {
	consumerService ConsumerService
	kafkaBrokers    string
	kafkaGroup      string
	kafkaTopic      string
	kafkaVersion    string
	kafkaOffset     string
}

func NewKafkaUtil(consumerService ConsumerService, kafkaBrokers, kafkaGroup, kafkaTopic, kafkaVersion, kafkaOffset string) KafkaUtil {
	return kafkaUtil{
		consumerService: consumerService,
		kafkaBrokers:    kafkaBrokers,
		kafkaGroup:      kafkaGroup,
		kafkaTopic:      kafkaTopic,
		kafkaVersion:    kafkaVersion,
		kafkaOffset:     kafkaOffset,
	}
}
func (k kafkaUtil) StartKafka(retryDelay int) {
	config, err := k.SetupConfig()
	if err != nil {
		return
	}

	for {
		var kafkaClient sarama.ConsumerGroup
		consumeError := make(chan error)

		for {
			kafkaClient, err = k.setupKafkaClient(config)
			if err == nil {
				break
			}

			time.Sleep((k.getDelay(retryDelay)) * time.Second)
		}

		k.startConsumerRoutine(kafkaClient, consumeError)

		consumerError := <-consumeError

		log.Println(consumerError)
	}

}

func (k kafkaUtil) getDelay(retryDelay int) time.Duration {
	if retryDelay < 1 {
		retryDelay = 3
	}
	return time.Duration(retryDelay)
}

func (k kafkaUtil) setupKafkaClient(config *sarama.Config) (sarama.ConsumerGroup, error) {
	brokersList := strings.Split(k.kafkaBrokers, ",")
	var client sarama.ConsumerGroup
	client, err := sarama.NewConsumerGroup(brokersList, k.kafkaGroup, config)
	if err != nil {
		return client, err
	}
	return client, nil
}

func (k kafkaUtil) startConsumerRoutine(client sarama.ConsumerGroup, consumeError chan error) {
	consumer := NewClientGroupMessageHandler(make(chan bool, 0), k.consumerService, k.kafkaTopic)
	go k.startConsumerWithWait(client, consumer, consumeError)
}

func (k kafkaUtil) startConsumerWithWait(client sarama.ConsumerGroup, consumer ClientGroupMessageHandler, consumeError chan error) {
	var err error
	ctx := context.TODO()
	for {
		err = client.Consume(ctx, strings.Split(k.kafkaTopic, ","), &consumer)
		if err != nil {
			break
		}
		if ctx.Err() != nil {
			err = ctx.Err()
			break
		}
		consumer.ready = make(chan bool)
	}
	consumeError <- err
}

func (k kafkaUtil) SetupConfig() (*sarama.Config, error) {
	config := sarama.NewConfig()
	version, err := validateKafkaVersion(k.kafkaVersion)
	if err != nil {
		return config, err
	}
	config.Version = version
	offSet, err := k.GetKafkaOffSet()
	if err != nil {
		return config, err
	}
	config.Consumer.Offsets.Initial = offSet
	return config, nil
}

func (k kafkaUtil) GetKafkaOffSet() (int64, error) {
	switch k.kafkaOffset {
	case OffsetOldest:
		return sarama.OffsetOldest, nil
	case OffsetNewest, "":
		return sarama.OffsetNewest, nil
	default:
		return 0, errors.New("Error in parsing offset value=" + k.kafkaOffset)

	}
}
func validateKafkaVersion(kafkaVersion string) (sarama.KafkaVersion, error) {
	var version sarama.KafkaVersion
	version, err := sarama.ParseKafkaVersion(kafkaVersion)
	if err != nil {
		return version, err
	}
	return version, nil
}
