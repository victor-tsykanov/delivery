package in

import "context"

type IMoveCouriersCommandHandler interface {
	Handle(ctx context.Context) error
}
