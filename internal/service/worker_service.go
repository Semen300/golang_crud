package service

import (
	"crud-go/internal/model"
	"crud-go/internal/repository"
	"fmt"
)

type IWorkerService interface {
	GetAllOrders(string, int) ([]model.Order, error)
	GetOrderById(string, int, int) (model.Order, error)
	SetTaskCompleted(string, int, int) error
}

type WorkerService struct {
	OrderRepository repository.IOrderRepository
	TaskRepository  repository.ITaskRepository
}

func NewWorkerService(or repository.IOrderRepository, tr repository.ITaskRepository) WorkerService {
	return WorkerService{OrderRepository: or, TaskRepository: tr}
}

func (ws WorkerService) GetAllOrders(login string, role int) ([]model.Order, error) {
	if role != 2 {
		return nil,
			fmt.Errorf("You are not authorized for this operation")
	}
	orders, ordersErr := ws.OrderRepository.GetOrdersByWorker(login)
	if ordersErr != nil {
		return nil,
			fmt.Errorf("Error getting orders by Worker: \n%w", ordersErr)
	}

	tasks, tasksErr := ws.TaskRepository.GetAllTasks()
	if tasksErr != nil {
		return nil,
			fmt.Errorf("Error getting tasks by Worker:\n%w", tasksErr)
	}

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

func (ws WorkerService) GetOrderById(login string, role int, id int) (model.Order, error) {
	if role != 2 {
		return model.Order{}, fmt.Errorf("You are not authorized for this operation")
	}

	order, ordersErr := ws.OrderRepository.GetOrderById(id)
	if ordersErr != nil {
		return model.Order{},
			fmt.Errorf("Error getting orders by Worker: \n%w", ordersErr)
	}

	tasks, tasksErr := ws.TaskRepository.GetTasksByContract(id)
	if tasksErr != nil {
		return model.Order{},
			fmt.Errorf("Error getting tasks by Worker:\n%w", tasksErr)
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

func (ws WorkerService) SetTaskCompleted(login string, role int, id int) error {
	if role != 2 {
		return fmt.Errorf("You are not authorized for this operation")
	}

	task, taskErr := ws.TaskRepository.GetTaskById(id)
	if taskErr != nil {
		return fmt.Errorf("Error getting task: \n%w", taskErr)
	}

	order, orderErr := ws.OrderRepository.GetOrderById(task.OrderID)
	if orderErr != nil {
		return fmt.Errorf("Error getting order: \n%w", orderErr)
	}

	tasks, tasksErr := ws.TaskRepository.GetTasksByContract(order.ID)
	if tasksErr != nil {
		return fmt.Errorf("Error getting task: \n%w", tasksErr)
	}

	var taskToChange *model.Task
	for _, tsk := range tasks {
		if tsk.Id == id {
			taskToChange = &tsk
		}
	}

	taskToChange.Finished = true
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
	if priceUnfinished == 0 {
		order.Status = 2
	}

	_, orderSaveErr := ws.OrderRepository.SaveOrder(order)
	if orderSaveErr != nil {
		return fmt.Errorf("Error updating order:\n%w", orderSaveErr)
	}

	_, taskSaveErr := ws.TaskRepository.SaveTask(*taskToChange)
	if taskSaveErr != nil {
		return fmt.Errorf("Error updating task:\n%w", taskSaveErr)
	}

	return nil
}
