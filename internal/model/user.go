package model

import "fmt"

type User struct {
	Login    string // From DB
	Password string // From DB
	Fio      string // From DB
}

func (u User) ToString() string {
	return fmt.Sprintf("User {Login: %s, Password: %s, Fio: %s}",
		u.Login,
		u.Password,
		u.Fio)
}
