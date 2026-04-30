package service

import "crud-go/internal/model"

type IWorkerService interface {
	GetAllOrders(string, int) ([]model.Order, error)
	GetOrderByID(string, int, int) (model.Order, error)
	SetTaskCompleted(string, int, int) error
}
