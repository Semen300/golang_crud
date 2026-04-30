package model

import "time"

type OrderCreationDTO struct {
	Deadline      time.Time         `json:"deadline" binding:"required"`
	CustomerLogin string            `json:"customerLogin" binding:"required"`
	Tasks         []TaskCreationDTO `json:"tasks" binding:"required"`
}
