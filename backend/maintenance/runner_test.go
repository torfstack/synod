package maintenance

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRunJobRunsImmediatelyAndPeriodically(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	completed := make(chan struct{}, 2)
	go runJob(ctx, Job{Name: "test", Interval: time.Millisecond, Run: func(context.Context) error {
		completed <- struct{}{}
		return nil
	}})

	for range 2 {
		select {
		case <-completed:
		case <-time.After(time.Second):
			t.Fatal("maintenance job did not run")
		}
	}
	cancel()
}

func TestRunJobDoesNotOverlapItself(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	started := make(chan struct{})
	release := make(chan struct{})
	var calls atomic.Int32
	go runJob(ctx, Job{Name: "test", Interval: time.Millisecond, Run: func(context.Context) error {
		calls.Add(1)
		close(started)
		<-release
		return nil
	}})

	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("maintenance job did not start")
	}
	time.Sleep(10 * time.Millisecond)
	require.Equal(t, int32(1), calls.Load())
	close(release)
}
