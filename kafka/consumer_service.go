package kafka

import (
	"context"
	"time"
)

type ConsumerService interface {
	ProcessIncomingMessage(context context.Context, rawMessage []byte, incomingTime time.Time) error
}
