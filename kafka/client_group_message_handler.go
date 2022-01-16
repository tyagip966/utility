package kafka

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"log"
	"runtime/debug"
	"time"
)

type ClientGroupMessageHandler struct {
	ready     chan bool
	service   ConsumerService
	topicName string
}

func (consumer *ClientGroupMessageHandler) WaitForStart() {
	<-consumer.ready
}

func (consumer *ClientGroupMessageHandler) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.ready)
	return nil
}

func (consumer *ClientGroupMessageHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *ClientGroupMessageHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		messageTime := message.Timestamp
		messageText := message.Value
		switch message.Topic {
		case consumer.topicName:
			err := consumer.processMessage(messageText, message.Topic, messageTime)
			if err != nil {
				log.Printf("Failed to process kafka message, re-queuing for further processing %+v \n", err)
				return nil
			}
		default:
			log.Printf("Not processing messages from unfollowed topics  %+v \n", messageText)
		}
		session.MarkMessage(message, "")
	}
	return nil
}

func (consumer *ClientGroupMessageHandler) processMessage(message []byte, topicName string, messageTime time.Time) error {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic while processing message %+v, \n %s", err, string(debug.Stack()))
		}
	}()

	ctx := consumer.configureTracingForKafkaConsumer(context.Background(), message)
	return consumer.service.ProcessIncomingMessage(ctx, message, messageTime)
}

func (consumer *ClientGroupMessageHandler) configureTracingForKafkaConsumer(ctx context.Context, message []byte) context.Context {
	var baseEvent map[string]interface{}
	err := json.Unmarshal(message, &baseEvent)
	if err != nil {

		return ctx
	}

	return ctx
}

func NewClientGroupMessageHandler(ready chan bool, service ConsumerService, topicName string) ClientGroupMessageHandler {
	return ClientGroupMessageHandler{
		ready:     ready,
		service:   service,
		topicName: topicName,
	}

}
