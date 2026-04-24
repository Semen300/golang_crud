package model

import "time"

type Order struct {
	ID                  int       // From DB
	Name                string    // From DB
	Deadline            time.Time // From DB
	ManagerLogin        string    // From DB
	WorkerLogin         string    // From DB
	CustomerLogin       string    // From DB
	PercentOfComplition float64   // Calculated in service
	PriseTotal          float64   // From DB
	PriceUnfinished     float64   // Calculated in service
	Status              int       // From DB
	Tasks               []Task    // calculated in service
}
