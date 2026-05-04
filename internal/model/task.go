package model

import "fmt"

type Task struct {
	Id       int    // From DB
	Name     string // From DB
	OrderID  int    // From DB
	ItemID   int    // From DB
	Amount   int    // From DB
	Finished bool   // From DB
	Price    int    // From DB
}

func (t Task) ToString() string {
	return fmt.Sprintf("Task {Id: %d, Name: %s, OrderID: %d, ItemID: %d, Amount: %d, Finished: %v, Price: %d}",
		t.Id,
		t.Name,
		t.OrderID,
		t.ItemID,
		t.Amount,
		t.Finished,
		t.Price)
}
