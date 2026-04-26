package model

import "fmt"

type Worker struct {
	User
	SuperiorLogin     string // From DB
	NumberOfContracts uint   // Calculated in service
}

func (w Worker) ToString() string {
	return fmt.Sprintf("Worker {Login: %s, Password: %s, Fio: %s, SuperiorLogin: %s, NumberOfContracts: %d}",
		w.Login,
		w.Password,
		w.Fio, w.SuperiorLogin,
		w.NumberOfContracts)
}
