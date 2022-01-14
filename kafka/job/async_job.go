package job

import (
	"context"
	"github.com/tyagip966/utility/fraazoError"
	"time"
)

type asyncJob struct {
	Job
}

func NewAsyncJob(job Job) Job {
	return asyncJob{job}
}

func (asyncJob asyncJob) Process(ctx context.Context, rawMessage []byte, incomingTime time.Time) (error *fraazoError.Error) {

	go func() {
		error = asyncJob.Job.Process(ctx, rawMessage, incomingTime)

		if error != nil {

		} else {

		}
	}()
	return nil
}
