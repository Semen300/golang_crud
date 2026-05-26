package service

import (
	"crud-go/internal/model"
	"crud-go/internal/repository"
	"fmt"
	"time"
)

type ICustomerService interface {
	GetOrdersByCustomer(login string, role int) ([]model.Order, error)
	GetOrderById(login string, role int, id int) (model.Order, error)
	CreateOrder(login string, role int, deadline time.Time) (int, error)
	DeleteOrder(login string, role int, id int) (int, error)
	GetItems(login string, role int) ([]model.Item, error)
	GetBasket(login string, role int) ([]model.TaskCreationDTO, error)
	SaveToBasket(login string, role int, item model.TaskCreationDTO) error
	DeleteFromBasket(login string, role int, id int) error
	ClearBasket(login string, role int) error
}

type CustomerService struct {
	OrderRepository  *repository.IOrderRepository
	TaskRepository   *repository.ITaskRepository
	BasketRepository *repository.IBasketRepository
	ItemRepository   *repository.IItemRepository
}

func NewCustomerService(or repository.IOrderRepository, tr repository.ITaskRepository, br repository.IBasketRepository, ir repository.IItemRepository) CustomerService {
	return CustomerService{OrderRepository: &or, TaskRepository: &tr, BasketRepository: &br, ItemRepository: &ir}
}

func (cs CustomerService) GetOrdersByCustomer(login string, role int) ([]model.Order, error) {
	if role != 1 {
		return nil, fmt.Errorf("You are not authorized for this operation")
	}
	or := *(cs.OrderRepository)
	tr := *(cs.TaskRepository)

	orders, orderRepoErr := or.GetOrdersByCustomer(login)
	if orderRepoErr != nil {
		return nil, fmt.Errorf("Error getting orders by customer login: \n%w", orderRepoErr)
	}

	tasks, taskRepoErr := tr.GetAllTasks()
	if taskRepoErr != nil {
		return nil, fmt.Errorf("Error getting tasks from TaskRepository: \n%w", taskRepoErr)
	}
	tasksByOrder := make(map[int][]model.Task)
	for _, task := range tasks {
		tasksByOrder[task.OrderID] = append(tasksByOrder[task.OrderID], task)
	}

	for i, order := range orders {
		orders[i].Tasks = tasksByOrder[order.ID]
		var numberItemsCompleted, numberItemsAll int
		var priceUnfinished int
		for _, task := range tasksByOrder[order.ID] {
			if task.Finished {
				numberItemsCompleted += task.Amount
			} else {
				priceUnfinished += task.Price
			}
			numberItemsAll += task.Amount
		}
		orders[i].PercentOfComplition = float64(numberItemsCompleted) / float64(numberItemsAll)
		orders[i].PriceUnfinished = priceUnfinished
	}
	return orders, nil
}

func (cs CustomerService) GetOrderById(login string, role int, id int) (model.Order, error) {
	if role != 1 {
		return model.Order{}, fmt.Errorf("You are not authorized for this operation")
	}
	or := *(cs.OrderRepository)
	tr := *(cs.TaskRepository)

	order, orderRepoErr := or.GetOrderById(id)
	if orderRepoErr != nil {
		return model.Order{}, fmt.Errorf("Error getting order by order ID: \n%w", orderRepoErr)
	}

	tasks, taskRepoErr := tr.GetTasksByContract(id)
	if taskRepoErr != nil {
		return model.Order{}, fmt.Errorf("Error getting tasks by order ID: \n%w", taskRepoErr)
	}
	var numberItemsCompleted, numberItemsAll int
	var priceUnfinished int
	for _, task := range tasks {
		if task.Finished {
			numberItemsCompleted += task.Amount
		} else {
			priceUnfinished += task.Price
		}
		numberItemsAll += task.Amount
	}

	order.PercentOfComplition = float64(numberItemsCompleted) / float64(numberItemsAll)
	order.PriceUnfinished = priceUnfinished
	order.Tasks = tasks
	return order, nil
}

func (cs CustomerService) CreateOrder(login string, role int, deadline time.Time) (int, error) {
	if role != 1 {
		return 0, fmt.Errorf("You are not authorized for this operation")
	}
	or := *(cs.OrderRepository)
	tr := *(cs.TaskRepository)
	br := *(cs.BasketRepository)

	price := 0
	basket, basketErr := br.GetBasket(login)
	if basketErr != nil {
		return 0, basketErr
	}
	orderDTO := model.OrderCreationDTO{
		Deadline: deadline,
		Tasks:    basket,
	}
	for _, item := range orderDTO.Tasks {
		price += item.ItemPrice * item.Amount
	}
	genericName := fmt.Sprintf("Заказ для пользователя %s до %02d.%02d.%4d на сумму %d.%d", login, orderDTO.Deadline.Day(), orderDTO.Deadline.Month(), orderDTO.Deadline.Year(), price/100, price%100)
	orderToSave := model.Order{
		Name:          genericName,
		Deadline:      orderDTO.Deadline,
		CustomerLogin: login,
		PriceTotal:    price,
		Status:        0,
	}

	tx, txErr := tr.Begin()
	if txErr != nil {
		return 0,
			fmt.Errorf("Error starting transaction: %w", txErr)
	}
	createdId, orderRepoErr := or.SaveOrder(orderToSave)
	if orderRepoErr != nil {
		tx.Rollback()
		return 0, fmt.Errorf("Error during saving order: \n%w", orderRepoErr)
	}

	for _, task := range orderDTO.Tasks {
		taskGenericName := fmt.Sprintf("Задача: изготовление %s до %02d.%02d.%4d", task.Name, orderDTO.Deadline.Day(), orderDTO.Deadline.Month(), orderDTO.Deadline.Year())
		taskToSave := model.Task{
			Name:     taskGenericName,
			OrderID:  createdId,
			ItemID:   task.ItemID,
			Amount:   task.Amount,
			Finished: false,
			Price:    task.Amount * task.ItemPrice,
		}
		_, taskRepoErr := tr.SaveTask(taskToSave)
		if taskRepoErr != nil {
			tx.Rollback()
			return 0, fmt.Errorf("Error during saving task %s: \n%w", taskToSave.ToString(), orderRepoErr)
		}
	}

	tx.Commit()
	return createdId, nil
}

func (cs CustomerService) DeleteOrder(login string, role int, id int) (int, error) {
	if role != 1 {
		return 0, fmt.Errorf("You are not authorized for this operation")
	}
	or := *(cs.OrderRepository)
	tr := *(cs.TaskRepository)

	var priceUnfinished int
	tasks, taskRepoErr := tr.GetTasksByContract(id)
	if taskRepoErr != nil {
		return 0, fmt.Errorf("Error getting tasks by order ID: \n%w", taskRepoErr)
	}

	tx, txErr := tr.Begin()
	if txErr != nil {
		return 0, fmt.Errorf("Error starting transaction")
	}
	for _, task := range tasks {
		if !task.Finished {
			priceUnfinished += task.Price
		}
		taskDeleteErr := tr.DeleteTask(task.ID)
		if taskDeleteErr != nil {
			tx.Rollback()
			return 0, fmt.Errorf("Error deleting task %s by order ID: \n%w", task.ToString(), taskDeleteErr)
		}
	}
	orderDeleteErr := or.DeleteOrder(id)
	if orderDeleteErr != nil {
		tx.Rollback()
		return 0, fmt.Errorf("Error deleting order by ID: \n%w", orderDeleteErr)
	}
	tx.Commit()
	return priceUnfinished, nil
}

func (cs CustomerService) GetItems(login string, role int) ([]model.Item, error) {
	if role != 1 {
		return nil, fmt.Errorf("You are not authorized for this operation")
	}
	ir := *(cs.ItemRepository)

	items, itemRepoErr := ir.GetAllItems()
	if itemRepoErr != nil {
		return nil, fmt.Errorf("Error getting items: \n%w", itemRepoErr)
	}

	return items, nil
}

func (cs CustomerService) GetBasket(login string, role int) ([]model.TaskCreationDTO, error) {
	if role != 1 {
		return nil,
			fmt.Errorf("You are not authorized for this operation")
	}
	br := *(cs.BasketRepository)

	basket, basketErr := br.GetBasket(login)
	if basketErr != nil {
		return nil,
			fmt.Errorf("Error getting customer basket: \n%w", basketErr)
	}

	return basket, nil
}

func (cs CustomerService) SaveToBasket(login string, role int, item model.TaskCreationDTO) error {
	if role != 1 {
		return fmt.Errorf("You are not authorized for this operation")
	}
	br := *(cs.BasketRepository)

	saveErr := br.SaveToBasket(login, item)
	if saveErr != nil {
		return fmt.Errorf("Error saving item to basket: \n%w", saveErr)
	}
	return nil
}

func (cs CustomerService) DeleteFromBasket(login string, role int, itemId int) error {
	if role != 1 {
		return fmt.Errorf("You are not authorized for this operation")
	}
	br := *(cs.BasketRepository)

	deleteErr := br.DeleteFromBasket(login, itemId)
	if deleteErr != nil {
		return fmt.Errorf("Error deleting item from basket: \n%w", deleteErr)
	}

	return nil
}

func (cs CustomerService) ClearBasket(login string, role int) error {
	if role != 1 {
		return fmt.Errorf("You are not authorized for this operation")
	}
	br := *cs.BasketRepository
	return br.ClearBasket(login)
}
