package service_test

import (
	"crud-go/internal/model"
	"crud-go/internal/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewWorkerService(t *testing.T) {
	orderRepoMock := new(orderRepoMock)
	taskRepoMock := new(taskRepoMock)
	testService := service.NewWorkerService(orderRepoMock, taskRepoMock)

	assert.NotEqual(t, service.WorkerService{}, testService)
}

func TestGetAllOrdersWorker(t *testing.T) {
	orderRepoMock := new(orderRepoMock)
	taskRepoMock := new(taskRepoMock)
	testService := service.NewWorkerService(orderRepoMock, taskRepoMock)

	mockOrders := []model.Order{
		{
			ID:            1,
			Name:          "order1",
			Deadline:      time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC),
			ManagerLogin:  "manager1",
			WorkerLogin:   "worker1",
			CustomerLogin: "customer1",
			PriceTotal:    10000,
			Status:        1,
		},
		{
			ID:            2,
			Name:          "order2",
			Deadline:      time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC),
			ManagerLogin:  "manager1",
			WorkerLogin:   "worker1",
			CustomerLogin: "customer2",
			PriceTotal:    20000,
			Status:        1,
		},
	}
	mockTasks := []model.Task{
		{
			ID:       1,
			Name:     "task1",
			OrderID:  1,
			ItemID:   1,
			Amount:   10,
			Finished: true,
			Price:    5000,
		},
		{
			ID:       2,
			Name:     "task2",
			OrderID:  1,
			ItemID:   2,
			Amount:   10,
			Finished: false,
			Price:    5000,
		},
		{
			ID:       3,
			Name:     "task3",
			OrderID:  2,
			ItemID:   3,
			Amount:   10,
			Finished: true,
			Price:    10000,
		},
		{
			ID:       4,
			Name:     "task4",
			OrderID:  2,
			ItemID:   4,
			Amount:   30,
			Finished: false,
			Price:    10000,
		},
	}

	orderRepoMock.On("GetOrdersByWorker", "worker1").Return(mockOrders, nil)
	taskRepoMock.On("GetAllTasks").Return(mockTasks, nil)

	expectedOrders := []model.Order{
		{
			ID:                  1,
			Name:                "order1",
			Deadline:            time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC),
			ManagerLogin:        "manager1",
			WorkerLogin:         "worker1",
			CustomerLogin:       "customer1",
			PercentOfComplition: 0.5,
			PriceTotal:          10000,
			PriceUnfinished:     5000,
			Status:              1,
			Tasks:               []model.Task{mockTasks[0], mockTasks[1]},
		},
		{
			ID:                  2,
			Name:                "order2",
			Deadline:            time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC),
			ManagerLogin:        "manager1",
			WorkerLogin:         "worker1",
			CustomerLogin:       "customer2",
			PercentOfComplition: 0.25,
			PriceTotal:          20000,
			PriceUnfinished:     10000,
			Status:              1,
			Tasks:               []model.Task{mockTasks[2], mockTasks[3]},
		},
	}

	gotOrders, gotErr := testService.GetAllOrders("worker1", 2)
	assert.NoError(t, gotErr)
	assert.Equal(t, expectedOrders, gotOrders)
}

func TestGetOrderByIdWorker(t *testing.T) {
	orderRepoMock := new(orderRepoMock)
	taskRepoMock := new(taskRepoMock)
	testService := service.NewWorkerService(orderRepoMock, taskRepoMock)

	mockOrder := model.Order{
		ID:            1,
		Name:          "order1",
		Deadline:      time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC),
		ManagerLogin:  "manager1",
		WorkerLogin:   "worker1",
		CustomerLogin: "customer1",
		PriceTotal:    10000,
		Status:        1,
	}

	mockTasks := []model.Task{
		{
			ID:       1,
			Name:     "task1",
			OrderID:  1,
			ItemID:   1,
			Amount:   10,
			Finished: true,
			Price:    1000,
		},
		{
			ID:       2,
			Name:     "task2",
			OrderID:  1,
			ItemID:   2,
			Amount:   90,
			Finished: false,
			Price:    9000,
		},
	}

	orderRepoMock.On("GetOrderById", 1).Return(mockOrder, nil)
	taskRepoMock.On("GetTasksByContract", 1).Return(mockTasks, nil)

	expectedOrder := model.Order{
		ID:                  1,
		Name:                "order1",
		Deadline:            time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC),
		ManagerLogin:        "manager1",
		WorkerLogin:         "worker1",
		CustomerLogin:       "customer1",
		PercentOfComplition: 0.1,
		PriceTotal:          10000,
		PriceUnfinished:     9000,
		Status:              1,
		Tasks:               mockTasks,
	}

	gotOrder, gotErr := testService.GetOrderById("worker1", 2, 1)
	assert.NoError(t, gotErr)
	assert.Equal(t, expectedOrder, gotOrder)
}

func TestSetTaskCompleted(t *testing.T) {
	orderRepoMock := new(orderRepoMock)
	taskRepoMock := new(taskRepoMock)
	testService := service.NewWorkerService(orderRepoMock, taskRepoMock)

	mockOrder := model.Order{
		ID:            1,
		Name:          "order1",
		Deadline:      time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC),
		ManagerLogin:  "manager1",
		WorkerLogin:   "worker1",
		CustomerLogin: "customer1",
		PriceTotal:    10000,
		Status:        1,
	}

	mockTasks := []model.Task{
		{
			ID:       1,
			Name:     "task1",
			OrderID:  1,
			ItemID:   1,
			Amount:   10,
			Finished: true,
			Price:    1000,
		},
		{
			ID:       2,
			Name:     "task2",
			OrderID:  1,
			ItemID:   2,
			Amount:   90,
			Finished: false,
			Price:    9000,
		},
	}

	changedTask := model.Task{
		ID:       2,
		Name:     "task2",
		OrderID:  1,
		ItemID:   2,
		Amount:   90,
		Finished: true,
		Price:    9000,
	}

	orderToSave := model.Order{
		ID:                  1,
		Name:                "order1",
		Deadline:            time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC),
		ManagerLogin:        "manager1",
		WorkerLogin:         "worker1",
		CustomerLogin:       "customer1",
		PercentOfComplition: 1,
		PriceTotal:          10000,
		PriceUnfinished:     0,
		Status:              2,
	}

	orderRepoMock.On("GetOrderById", 1).Return(mockOrder, nil)
	orderRepoMock.On("SaveOrder", orderToSave).Return(1, nil)
	taskRepoMock.On("GetTaskById", 2).Return(mockTasks[1], nil)
	taskRepoMock.On("GetTasksByContract", 1).Return(mockTasks, nil)
	taskRepoMock.On("SaveTask", changedTask).Return(1, nil)

	gotErr := testService.SetTaskCompleted("worker1", 2, 2)
	assert.NoError(t, gotErr)
}
