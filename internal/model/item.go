package model

import "fmt"

type Item struct {
	ID    int    // From DB
	Name  string // From DB
	Price uint   // From DB
}

func (i Item) ToString() string {
	return fmt.Sprintf("Item {Id: %d, Name: %s, Price: %d}",
		i.ID,
		i.Name,
		i.Price)
}
