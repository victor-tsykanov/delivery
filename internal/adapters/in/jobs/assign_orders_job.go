package jobs

import (
	"context"

	inPorts "github.com/victor-tsykanov/delivery/internal/core/ports/in"
)

type AssignOrdersJob struct {
	assignOrdersCommandHandler inPorts.IAssignOrdersCommandHandler
}

func NewAssignOrdersJob(assignOrdersCommandHandler inPorts.IAssignOrdersCommandHandler) (*AssignOrdersJob, error) {
	return &AssignOrdersJob{assignOrdersCommandHandler: assignOrdersCommandHandler}, nil
}

func (j *AssignOrdersJob) Execute(ctx context.Context) error {
	return j.assignOrdersCommandHandler.Handle(ctx)
}
