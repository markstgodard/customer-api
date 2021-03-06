package customers

import (
	"context"
	"errors"
)

var (
	ErrInvalid = errors.New("invalid request")
)

type Service interface {
	CreateCustomer(ctx context.Context, c Customer) (Customer, error)
	ListCustomers(ctx context.Context) ([]Customer, error)
}

type customerService struct {
	repo CustomerRepository
}

func NewService(repo CustomerRepository) Service {
	return &customerService{repo}
}

func (s *customerService) CreateCustomer(_ context.Context, c Customer) (Customer, error) {
	// validate customer
	if ok, _ := c.Valid(); !ok {
		return c, ErrInvalid
	}

	// create customer
	err := s.repo.Create(c)
	if err != nil {
		return c, err
	}

	return c, nil
}

func (s *customerService) ListCustomers(_ context.Context) ([]Customer, error) {
	return s.repo.List()
}
