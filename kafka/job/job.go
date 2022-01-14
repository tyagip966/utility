package job

import (
	"context"
	"fraazo-utility/fraazoError"
	"time"
)

type Job interface {
	GetName() string
	Process(ctx context.Context, rawMessage []byte, incomingTime time.Time) (error *fraazoError.Error)
}
