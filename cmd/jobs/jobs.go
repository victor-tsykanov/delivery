package jobs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/victor-tsykanov/delivery/cmd/app"
	"github.com/victor-tsykanov/delivery/internal/adapters/in/jobs"
)

func AssignOrders(ctx context.Context, root *app.CompositionRoot) {
	job, err := jobs.NewAssignOrdersJob(root.CommandHandlers.AssignOrdersCommandHandler)
	if err != nil {
		log.Fatalf("failed to create AssignOrdersJob: %v", err)
	}

	repeatTask(ctx, job.Execute, 2*time.Second)
}

func MoveCouriers(ctx context.Context, root *app.CompositionRoot) {
	job, err := jobs.NewMoveCouriersJob(root.CommandHandlers.MoveCouriersCommandHandler)
	if err != nil {
		log.Fatalf("failed to create MoveCouriersJob: %v", err)
	}

	repeatTask(ctx, job.Execute, 1*time.Second)
}

func repeatTask(ctx context.Context, fn func(ctx context.Context) error, interval time.Duration) {
	for {
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("recovering after task panic:", r)
					time.Sleep(interval)
				}
			}()

			err := fn(ctx)
			if err != nil {
				fmt.Println("task failed with error:", err)
			}

			time.Sleep(interval)
		}()
	}
}
