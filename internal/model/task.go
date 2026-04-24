package model

type Task struct {
	Id       int
	Name     string
	Item     *Item
	Amount   int
	Finished bool
	Price    uint
}
