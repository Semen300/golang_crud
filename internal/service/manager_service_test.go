package service_test

import (
	"crud-go/internal/model"
	"crud-go/internal/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewManagerService(t *testing.T) {
	orderRepoMock := new(orderRepoMock)
	taskRepoMock := new(taskRepoMock)
	userRepoMock := new(userRepoMock)
	testService := service.NewManagerService(orderRepoMock, taskRepoMock, userRepoMock)

	assert.NotEqual(t, service.ManagerService{}, testService)
}

func TestGetAllWorkers(t *testing.T) {
	userRepoMock := new(userRepoMock)
	testService := service.NewManagerService(nil, nil, userRepoMock)

	mockWorkers := []model.Worker{
		model.NewWorker("worker1", "1111", "AAAA", "manager1"),
		model.NewWorker("worker2", "2222", "BBBB", "manager1"),
	}
	userRepoMock.On("GetWorkersByManager", "manager1").Return(mockWorkers, nil)

	workers, err := testService.GetAllWorkers("manager1", 3)
	assert.NoError(t, err)
	assert.Equal(t, mockWorkers, workers)
}

func TestGetAllOrders(t *testing.T) {
	orderRepoMock := new(orderRepoMock)
	taskRepoMock := new(taskRepoMock)

	testService := service.NewManagerService(orderRepoMock, taskRepoMock, nil)

	mockOrdersManager := []model.Order{
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
			WorkerLogin:   "worker2",
			CustomerLogin: "customer2",
			PriceTotal:    20000,
			Status:        1,
		},
	}

	mockOrdersEmpty := []model.Order{
		{
			ID:            3,
			Name:          "order3",
			Deadline:      time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC),
			CustomerLogin: "customer3",
			PriceTotal:    30000,
			Status:        0,
		},
	}

	orderRepoMock.On("GetOrdersByManager", "manager1").Return(mockOrdersManager, nil)
	orderRepoMock.On("GetOrdersByManager", "").Return(mockOrdersEmpty, nil)

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
		{
			ID:       5,
			Name:     "task5",
			OrderID:  3,
			ItemID:   5,
			Amount:   500,
			Finished: false,
			Price:    30000,
		},
	}

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
			WorkerLogin:         "worker2",
			CustomerLogin:       "customer2",
			PercentOfComplition: 0.25,
			PriceTotal:          20000,
			PriceUnfinished:     10000,
			Status:              1,
			Tasks:               []model.Task{mockTasks[2], mockTasks[3]},
		},
		{
			ID:                  3,
			Name:                "order3",
			Deadline:            time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC),
			CustomerLogin:       "customer3",
			PercentOfComplition: 0,
			PriceTotal:          30000,
			PriceUnfinished:     30000,
			Status:              0,
			Tasks:               []model.Task{mockTasks[4]},
		},
	}

	ordersGot, err := testService.GetAllOrders("manager1", 3)
	assert.NoError(t, err)
	assert.Equal(t, expectedOrders, ordersGot)
}

func TestManagerGetOrderById(t *testing.T) {
	orderRepoMock := new(orderRepoMock)
	taskRepoMock := new(taskRepoMock)

	testService := service.NewManagerService(orderRepoMock, taskRepoMock, nil)

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

	orderRepoMock.On("GetOrderById", 1).Return(mockOrder, nil)

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

	orderGot, err := testService.GetOrderById("manager1", 3, 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedOrder, orderGot)

}

func TestAssignWorkerToOrder(t *testing.T) {
	orderRepoMock := new(orderRepoMock)
	testService := service.NewManagerService(orderRepoMock, nil, nil)

	mockOrder := model.Order{
		ID:           1,
		Name:         "order1",
		Deadline:     time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC),
		ManagerLogin: "manager1",
		PriceTotal:   10000,
		Status:       0,
	}

	mockOrderChanged := mockOrder
	mockOrderChanged.ManagerLogin = "manager1"
	mockOrderChanged.WorkerLogin = "worker2"
	mockOrderChanged.Status = 1

	orderRepoMock.On("GetOrderById", 1).Return(mockOrder, nil)
	orderRepoMock.On("SaveOrder", mockOrderChanged).Return(1, nil)

	err := testService.AssignWorkerToOrder("manager1", 3, 1, "worker2")
	assert.NoError(t, err)
}
