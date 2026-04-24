package model

type Worker struct {
	User
	SuperiorLogin     string // From DB
	NumberOfContracts uint   // Calculated in service
}
