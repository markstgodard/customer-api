package customers

import (
	"context"

	jujuratelimit "github.com/juju/ratelimit"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
)

var qps = 100

// "make" methods for each endpoint
func makeGetCustomersEndpoint(s Service) endpoint.Endpoint {
	e := getCustomersEndpoint(s)

	e = ratelimitingMiddleware(qps, e)
	return e
}

func makeCreateCustomerEndpoint(s Service) endpoint.Endpoint {
	e := createCustomerEndpoint(s)

	e = ratelimitingMiddleware(qps, e)
	return e
}

// rate limiting middleware
func ratelimitingMiddleware(qps int, e endpoint.Endpoint) endpoint.Endpoint {
	return ratelimit.NewTokenBucketLimiter(jujuratelimit.NewBucketWithRate(float64(qps), int64(qps)))(e)
}

// Get Customers
type getCustomersRequest struct {
}

type getCustomersResponse struct {
	Customers []Customer `json:"customers,omitempty"`
	Err       error      `json:"err,omitempty"`
}

func getCustomersEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		customerz, err := s.ListCustomers(ctx)
		if err != nil {
			return errResponse{err}, err
		}
		return getCustomersResponse{Customers: customerz, Err: err}, nil
	}
}

// Create Customer
type createCustomerRequest struct {
	Customer
}

func createCustomerEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(createCustomerRequest)
		c, err := s.CreateCustomer(ctx, req.Customer)
		if err != nil {
			return errResponse{err}, err
		}
		return c, nil
	}
}

type errResponse struct {
	Err error `json:"err,omitempty"`
}
