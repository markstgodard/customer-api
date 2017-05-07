package customers

import "fmt"

type CustomerRepository interface {
	Create(Customer) error
	List() ([]Customer, error)
}

type InMemoryCustomerRepo struct {
}

func (i *InMemoryCustomerRepo) Create(c Customer) error {
	// not really
	fmt.Printf("creating customer: %v\n", c)
	return nil
}

func (i *InMemoryCustomerRepo) List() ([]Customer, error) {
	// not implemented
	fmt.Println("find all customers")

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
