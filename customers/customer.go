package customers

type Customer struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func (c Customer) Valid() (bool, error) {
	// validate customer
	return true, nil
}
