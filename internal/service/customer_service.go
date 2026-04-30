package service

import "crud-go/internal/model"

type ICustomerService interface {
	GetOrdersByCustomer(login string, role int) ([]model.Order, error)
	GetOrderByID(login string, role int, id int) (model.Order, error)
	CreateOrder(login string, role int, orderDTO model.OrderCreationDTO) (int, error)
	DeleteOrder(login string, role int, id int) error
	GetItems(login string, role int) ([]model.Item, error)
	GetBasket(login string, role int) ([]model.Task, error)
	SaveToBasket(login string, role int, item model.TaskCreationDTO) error
	DeleteFromBasket(login string, role int, id int) error
}
