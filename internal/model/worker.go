package model

import "fmt"

type Worker struct {
	User
	SuperiorLogin string // From DB
}

func NewWorker(login, password, fio, superiorLogin string) Worker {
	return Worker{User{login, password, fio}, superiorLogin}
}

func (w Worker) ToString() string {
	return fmt.Sprintf("Worker {Login: %s, Password: %s, Fio: %s, SuperiorLogin: %s}",
		w.Login,
		w.Password,
		w.Fio, w.SuperiorLogin)
}
