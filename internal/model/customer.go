package model

import "fmt"

type Customer struct {
	User
	Number string // From DB
	Email  string // From DB
}

func (c Customer) ToString() string {
	return fmt.Sprintf("Customer {Login: %s, Password: %s, Fio: %s, Number: %s, Email: %s}",
		c.Login,
		c.Password,
		c.Fio,
		c.Number,
		c.Email)
}
