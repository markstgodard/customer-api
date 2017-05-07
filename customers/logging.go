package customers

import (
	"context"

	"github.com/go-kit/kit/log"
)

type loggingMiddleware struct {
	logger log.Logger
	next   Service
}

func NewLoggingMiddlware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return loggingMiddleware{logger, next}
	}
}

func (l loggingMiddleware) Create(ctx context.Context, c Customer) (v Customer, err error) {
	defer func() {
		l.logger.Log("method", "Create", "email", c.Email, "err", err)
	}()
	return l.next.Create(ctx, c)
}

func (l loggingMiddleware) List(ctx context.Context) (v []Customer, err error) {
	defer func() {
		l.logger.Log("method", "List", "err", err)
	}()
	return l.next.List(ctx)
}
