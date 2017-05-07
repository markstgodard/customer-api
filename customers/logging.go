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

func (l loggingMiddleware) CreateCustomer(ctx context.Context, c Customer) (v Customer, err error) {
	defer func() {
		l.logger.Log("method", "CreateCustomer", "email", c.Email, "err", err)
	}()
	return l.next.CreateCustomer(ctx, c)
}

func (l loggingMiddleware) ListCustomers(ctx context.Context) (v []Customer, err error) {
	defer func() {
		l.logger.Log("method", "ListCustomers", "err", err)
	}()
	return l.next.ListCustomers(ctx)
}
