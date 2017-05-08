package customers

import (
	"context"
	"time"

	"github.com/go-kit/kit/metrics"
)

type metricsService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           Service
}

func NewMetricsService(counter metrics.Counter, latency metrics.Histogram, s Service) Service {
	return &metricsService{
		requestCount:   counter,
		requestLatency: latency,
		next:           s,
	}
}

func (s *metricsService) CreateCustomer(ctx context.Context, c Customer) (v Customer, err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "createCustomer").Add(1)
		s.requestLatency.With("method", "createCustomer").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.next.CreateCustomer(ctx, c)
}

func (s *metricsService) ListCustomers(ctx context.Context) (v []Customer, err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "listCustomers").Add(1)
		s.requestLatency.With("method", "listCustomers").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.next.ListCustomers(ctx)
}
