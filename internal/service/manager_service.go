package service

import (
	"crud-go/internal/model"
	"crud-go/internal/repository"
	"fmt"
)

type IManagerService interface {
	GetAllWorkers(string, int) ([]model.Worker, error)
	GetAllOrders(string, int) ([]model.Order, error)
	GetOrderById(string, int, int) (model.Order, error)
	AssignWorkerToOrder(string, int, int, string) error
}

type ManagerService struct {
	OrderRepository repository.IOrderRepository
	TaskRepository  repository.ITaskRepository
	UserRepository  repository.IUserRepository
}

func NewManagerService(or repository.IOrderRepository, tr repository.ITaskRepository, ur repository.IUserRepository) ManagerService {
	return ManagerService{OrderRepository: or, TaskRepository: tr, UserRepository: ur}
}

func (ms ManagerService) GetAllWorkers(login string, role int) ([]model.Worker, error) {
	if role != 3 {
		return nil,
			fmt.Errorf("You are not authorized for this operation")
	}

	workers, getErr := ms.UserRepository.GetWorkersByManager(login)
	if getErr != nil {
		return nil,
			fmt.Errorf("Error getting workers by manager: \n%w", getErr)
	}

	return workers, nil
}

func (ms ManagerService) GetAllOrders(login string, role int) ([]model.Order, error) {
	if role != 3 {
		return nil, fmt.Errorf("You are not authorized for this operation")
	}

	managerOrders, ordersErr := ms.OrderRepository.GetOrdersByManager(login)
	emptyOrders, ordersErr := ms.OrderRepository.GetOrdersByManager("")
	if ordersErr != nil {
		return nil,
			fmt.Errorf("Error getting orders for manager %s:\n%w", login, ordersErr)
	}

	tasks, tasksErr := ms.TaskRepository.GetAllTasks()
	if tasksErr != nil {
		return nil,
			fmt.Errorf("Error getting tasks:\n%w", tasksErr)
	}

	orders := append(managerOrders, emptyOrders...)
	tasksMap := make(map[int][]model.Task, 0)
	for _, task := range tasks {
		tasksMap[task.OrderID] = append(tasksMap[task.OrderID], task)
	}

	for i, order := range orders {
		orders[i].Tasks = tasksMap[order.ID]
		var numberItemsCompleted, numberItemsAll int
		var priceUnfinished int
		for _, task := range tasksMap[order.ID] {
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

func (ms ManagerService) GetOrderById(login string, role int, id int) (model.Order, error) {
	if role != 3 {
		return model.Order{}, fmt.Errorf("You are not authorized for this operation")
	}

	order, orderErr := ms.OrderRepository.GetOrderById(id)
	if orderErr != nil {
		return model.Order{},
			fmt.Errorf("Error getting order by ID: \n%w", orderErr)
	}
	tasks, tasksErr := ms.TaskRepository.GetTasksByContract(id)
	if tasksErr != nil {
		return model.Order{},
			fmt.Errorf("Error getting tasks for order: \n%w", tasksErr)
	}
	order.Tasks = tasks
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

	return order, nil
}

func (ms ManagerService) AssignWorkerToOrder(login string, role int, id int, workerLogin string) error {
	if role != 3 {
		return fmt.Errorf("You are not authorized for this operation")
	}
	order, getErr := ms.OrderRepository.GetOrderById(id)
	if getErr != nil {
		return fmt.Errorf("Error getting order by ID: \n%w", getErr)
	}
	order.ManagerLogin = login
	order.WorkerLogin = workerLogin
	order.Status = 1
	_, saveErr := ms.OrderRepository.SaveOrder(order)
	if saveErr != nil {
		return fmt.Errorf("Error saving order: \n%w", saveErr)
	}
	return nil
}
