package ports

import "context"

type ITransactionManager interface {
	Execute(ctx context.Context, fn func(context.Context) error) error
}
