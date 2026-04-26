package model

import "fmt"

type Manager struct {
	User
}

func (m Manager) ToString() string {
	return fmt.Sprintf("Manager {Login: %s, Password: %s, Fio: %s}",
		m.Login,
		m.Password,
		m.Fio)
}
