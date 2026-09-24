package maintenance

import (
	"context"
	"time"

	"github.com/torfstack/synod/backend/logging"
)

type Job struct {
	Name     string
	Interval time.Duration
	Run      func(context.Context) error
}

func Run(ctx context.Context, jobs ...Job) {
	for _, job := range jobs {
		go runJob(ctx, job)
	}
	<-ctx.Done()
}

func runJob(ctx context.Context, job Job) {
	run := func() {
		if err := job.Run(ctx); err != nil {
			logging.Errorf(ctx, "maintenance job %s failed: %v", job.Name, err)
		}
	}
	run()
	ticker := time.NewTicker(job.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			run()
		}
	}
}
