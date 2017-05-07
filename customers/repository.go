package customers

type CustomerRepository interface {
	Create(Customer) error
	List() ([]Customer, error)
}

type InMemoryCustomerRepo struct {
}

func (i *InMemoryCustomerRepo) Create(c Customer) error {
	// not really
	return nil
}

func (i *InMemoryCustomerRepo) List() ([]Customer, error) {
	// not implemented

	c := []Customer{
		{
			ID:        1,
			FirstName: "Mark",
			LastName:  "St.Godard",
			Email:     "markstgodard@gmail.com",
		},
	}

	return c, nil
}
