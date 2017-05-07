package customers

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type errResponse struct {
	Err error `json:"err,omitempty"`
}

type getCustomersRequest struct {
}

type getCustomersResponse struct {
	Customers []Customer `json:"customers,omitempty"`
	Err       error      `json:"err,omitempty"`
}

func makeGetCustomersEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		customerz, err := s.List(ctx)
		if err != nil {
			return errResponse{err}, err
		}
		return getCustomersResponse{Customers: customerz, Err: err}, nil
	}
}

type createCustomerRequest struct {
	Customer
}

func makeCreateCustomerEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(createCustomerRequest)
		c, err := s.Create(ctx, req.Customer)
		if err != nil {
			return errResponse{err}, err
		}
		return c, nil
	}
}
