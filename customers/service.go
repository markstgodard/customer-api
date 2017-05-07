package customers

import "context"

type Service interface {
	Create(ctx context.Context, c Customer) (Customer, error)
	List(ctx context.Context) ([]Customer, error)
}

type customerService struct {
	repo CustomerRepository
}

func NewService(repo CustomerRepository) Service {
	return &customerService{repo}
}

func (s *customerService) Create(_ context.Context, c Customer) (Customer, error) {
	// validate customer
	if ok, err := c.Valid(); !ok {
		return c, err
	}

	// create customer
	err := s.repo.Create(c)
	if err != nil {
		return c, err
	}

	return c, nil
}

func (s *customerService) List(_ context.Context) ([]Customer, error) {
	return s.repo.List()
}
