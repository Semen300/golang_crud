package model

type Task struct {
	Id         int    // From DB
	Name       string // From DB
	ContractID int    // From DB
	ItemID     int    // From DB
	Amount     int    // From DB
	Finished   bool   // From DB
	Price      uint   // From DB
}
