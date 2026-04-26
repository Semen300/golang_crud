package model

import "fmt"

type Item struct {
	Id    int    // From DB
	Name  string // From DB
	Price uint   // From DB
}

func (i Item) ToString() string {
	return fmt.Sprintf("Item {Id: %d, Name: %s, Price: %d}",
		i.Id,
		i.Name,
		i.Price)
}
