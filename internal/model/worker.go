package model

type Worker struct {
	User
	Superior          Manager // From DB
	NumberOfContracts uint    // Calculated in service
}
