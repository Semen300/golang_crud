package model

import "time"

type OrderCreationDTO struct {
	Deadline time.Time         `json:"deadline" binding:"required"`
	Tasks    []TaskCreationDTO `json:"tasks" binding:"required"`
}
