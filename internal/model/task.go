package model

import "fmt"

type Task struct {
	Id         int    // From DB
	Name       string // From DB
	ContractID int    // From DB
	ItemID     int    // From DB
	Amount     int    // From DB
	Finished   bool   // From DB
	Price      uint   // From DB
}

func (t Task) ToString() string {
	return fmt.Sprintf("Task {Id: %d, Name: %s, ContractID: %d, ItemID: %d, Amount: %d, Finished: %v, Price: %d}",
		t.Id,
		t.Name,
		t.ContractID,
		t.ItemID,
		t.Amount,
		t.Finished,
		t.Price)
}
