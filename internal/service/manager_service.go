package service

import "crud-go/internal/model"

type IManagerService interface {
	GetAllOrders(string, int) ([]model.Order, error)
	GetOrderByID(string, int, int) (model.Order, error)
	SetWorkerLogin(string, int, int, string) error
}
