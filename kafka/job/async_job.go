package job

import (
	"context"
	"github.com/tyagip966/utility/fraazoError"
	"log"
	"time"
)

type asyncJob struct {
	Job
}

func NewAsyncJob(job Job) Job {
	return asyncJob{job}
}

func (asyncJob asyncJob) Process(ctx context.Context, rawMessage []byte, incomingTime time.Time) (error *fraazoError.Error) {
	log.Printf("Inside AsyncJob: Scheduling Job: %s \n", asyncJob.GetName())
	go func() {
		log.Printf("AsyncJob: Job %s kick starting \n", asyncJob.GetName())
		error = asyncJob.Job.Process(ctx, rawMessage, incomingTime)

		if error != nil {
			log.Printf("AsyncJob: Job %s completed with error. Error: %+v \n", asyncJob.GetName(), error)
		} else {
			log.Printf("AsyncJob: Job %s completed successfully \n", asyncJob.GetName())
		}
	}()
	return nil
}
