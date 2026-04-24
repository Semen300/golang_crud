package model

import "time"

type Contract struct {
	ID                  int       // From DB
	Name                string    // From DB
	Deadline            time.Time // From DB
	Manager             *Manager  // From DB JOIN
	Worker              *Worker   // From DB JOIN
	Customer            *Customer // From DB JOIN
	PercentOfComplition float64   // Calculated in service
	PriseTotal          float64   // From DB
	PriceUnfinished     float64   // Calculated in service
	Status              int       // From DB
	Tasks               []Task    // calculated in service
}
