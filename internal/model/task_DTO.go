package model

type TaskCreationDTO struct {
	ItemID int `json:"itemId" binding:"required"`
	Amount int `json:"amount" binding:"required"`
}
