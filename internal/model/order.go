package model

import (
	"fmt"
	"time"
)

type Order struct {
	ID                  int       // From DB
	Name                string    // From DB
	Deadline            time.Time // From DB
	ManagerLogin        string    // From DB
	WorkerLogin         string    // From DB
	CustomerLogin       string    // From DB
	PercentOfComplition float64   // Calculated in service
	PriseTotal          int       // From DB
	PriceUnfinished     int       // Calculated in service
	Status              int       // From DB
	Tasks               []Task    // calculated in service
}

func (o Order) ToString() string {
	return fmt.Sprintf("Order {ID: %d, Name: %s, Deadline: %v, ManagerLogin: %s, WorkerLogin: %s, CustomerLogin: %s, PercentOfconplition: %.2f, PriceTotal: %d, PriceUnfinished: %d, Status: %d, Tasks: %v}",
		o.ID,
		o.Name,
		o.Deadline,
		o.ManagerLogin,
		o.WorkerLogin,
		o.CustomerLogin,
		o.PercentOfComplition,
		o.PriseTotal,
		o.PriceUnfinished,
		o.Status,
		o.Tasks)
}
