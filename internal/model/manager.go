package model

import "fmt"

type Manager struct {
	User
}

func NewManager(login, password, fio string) Manager {
	return Manager{User{login, password, fio}}
}

func (m Manager) ToString() string {
	return fmt.Sprintf("Manager {Login: %s, Password: %s, Fio: %s}",
		m.Login,
		m.Password,
		m.Fio)
}
