package in

import "context"

type IAssignOrdersCommandHandler interface {
	Handle(context.Context) error
}
