package model

type Customer struct {
	User
	Number string // From DB
	Email  string // From DB
}
