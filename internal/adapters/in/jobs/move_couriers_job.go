package jobs

import (
	"context"

	inPorts "github.com/victor-tsykanov/delivery/internal/core/ports/in"
)

type MoveCouriersJob struct {
	moveCouriersCommandHandler inPorts.IMoveCouriersCommandHandler
}

func NewMoveCouriersJob(moveCouriersCommandHandler inPorts.IMoveCouriersCommandHandler) (*MoveCouriersJob, error) {
	return &MoveCouriersJob{moveCouriersCommandHandler: moveCouriersCommandHandler}, nil
}

func (j *MoveCouriersJob) Execute(ctx context.Context) error {
	return j.moveCouriersCommandHandler.Handle(ctx)
}
