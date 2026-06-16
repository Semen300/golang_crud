package model

type TaskCreationDTO struct {
	ItemID    int    `json:"itemId" binding:"required"`
	Name      string `json:"name" binding:"required"`
	ItemPrice int    `json:"itemPrice" binding:"required"`
	Amount    int    `json:"amount" binding:"required"`
}
