package job

import (
	"context"
	"time"
	"utility/fraazoError"
)

type Job interface {
	GetName() string
	Process(ctx context.Context, rawMessage []byte, incomingTime time.Time) (error *fraazoError.Error)
}
