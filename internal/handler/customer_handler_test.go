package handler_test

import (
	"crud-go/internal/handler"
	"crud-go/internal/model"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const customerPrefix = "/api/v1/customer"

type mockCustomerService struct {
	mock.Mock
}

func (m *mockCustomerService) GetOrdersByCustomer(login string, role int) ([]model.Order, error) {
	args := m.Called(login, role)
	return args.Get(0).([]model.Order), args.Error(1)
}

func (m *mockCustomerService) GetOrderByID(login string, role int, id int) (model.Order, error) {
	args := m.Called(login, role, id)
	return args.Get(0).(model.Order), args.Error(1)
}

func (m *mockCustomerService) CreateOrder(login string, role int, orderDTO model.OrderCreationDTO) (int, error) {
	args := m.Called(login, role, orderDTO)
	return args.Int(0), args.Error(1)
}

func (m *mockCustomerService) DeleteOrder(login string, role int, id int) error {
	args := m.Called(login, role, id)
	return args.Error(0)
}

func (m *mockCustomerService) GetItems(login string, role int) ([]model.Item, error) {
	args := m.Called(login, role)
	return args.Get(0).([]model.Item), args.Error(1)
}

func (m *mockCustomerService) GetBasket(login string, role int) ([]model.TaskCreationDTO, error) {
	args := m.Called(login, role)
	return args.Get(0).([]model.TaskCreationDTO), args.Error(1)
}

func (m *mockCustomerService) SaveToBasket(login string, role int, item model.TaskCreationDTO) error {
	args := m.Called(login, role, item)
	return args.Error(0)
}

func (m *mockCustomerService) DeleteFromBasket(login string, role int, id int) error {
	args := m.Called(login, role, id)
	return args.Error(0)
}

func TestGetAllOrders200(t *testing.T) {
	// Создаём мок-сервис и настраиваем его для возврата тестовых данных
	mockService := new(mockCustomerService)
	mockOrders := []model.Order{
		{ID: 1, Name: "Order 1", Deadline: time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC), ManagerLogin: "manager1", WorkerLogin: "worker1", CustomerLogin: "customer1", PercentOfComplition: 50.0, PriseTotal: 100, PriceUnfinished: 50, Status: 2},
		{ID: 2, Name: "Order 2", Deadline: time.Date(2026, time.May, 1, 0, 0, 0, 0, time.UTC), ManagerLogin: "manager2", WorkerLogin: "worker2", CustomerLogin: "customer1", PercentOfComplition: 100.0, PriseTotal: 200, PriceUnfinished: 0, Status: 3},
	}
	mockService.On("GetOrdersByCustomer", "customer1", 1).Return(mockOrders, nil)

	r := gin.New()
	r.GET(customerPrefix+"/orders", func(ctx *gin.Context) {
		ctx.Set("login", "customer1")
		ctx.Set("role", 1)
		handler.NewCustomerHandler(mockService).GetAllOrders(ctx)
	})

	req := httptest.NewRequest(http.MethodGet, customerPrefix+"/orders", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var responseOrders []model.Order
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &responseOrders))
	assert.Equal(t, len(mockOrders), len(responseOrders))
	assert.Equal(t, mockOrders, responseOrders)
}

func TestGet401(t *testing.T) {
	// Создаём пустой мок-сервис
	mockService := new(mockCustomerService)

	r := gin.New()
	r.GET(customerPrefix+"/orders", func(ctx *gin.Context) {
		handler.NewCustomerHandler(mockService).GetAllOrders(ctx)
	})

	req := httptest.NewRequest(http.MethodGet, customerPrefix+"/orders", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGet500(t *testing.T) {
	// Создаём пустой мок-сервис
	mockService := new(mockCustomerService)

	r := gin.New()
	r.GET(customerPrefix+"/orders", func(ctx *gin.Context) {
		ctx.Set("login", "customer1")
		ctx.Set("role", "super")
		handler.NewCustomerHandler(mockService).GetAllOrders(ctx)
	})

	req := httptest.NewRequest(http.MethodGet, customerPrefix+"/orders", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetCustomerOrder(t *testing.T) {
	// Создаём мок-сервис и настраиваем его для возврата тестовых данных
	mockService := new(mockCustomerService)
	mockOrder := model.Order{ID: 1, Name: "Order 1", Deadline: time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC), ManagerLogin: "manager1", WorkerLogin: "worker1", CustomerLogin: "customer1", PercentOfComplition: 50.0, PriseTotal: 100, PriceUnfinished: 50, Status: 2}
	mockService.On("GetOrderByID", "customer1", 1, 1).Return(mockOrder, nil)

	r := gin.New()
	r.GET(customerPrefix+"/orders/:id", func(ctx *gin.Context) {
		ctx.Set("login", "customer1")
		ctx.Set("role", 1)
		handler.NewCustomerHandler(mockService).GetCustomerOrder(ctx)
	})

	req := httptest.NewRequest(http.MethodGet, customerPrefix+"/orders/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var order model.Order
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &order))
	assert.Equal(t, mockOrder, order)
}

func TestGetCustomerOrderBadRequest(t *testing.T) {
	mockService := new(mockCustomerService)

	r := gin.New()
	r.GET(customerPrefix+"/orders/:id", func(ctx *gin.Context) {
		ctx.Set("login", "customer1")
		ctx.Set("role", 1)
		handler.NewCustomerHandler(mockService).GetCustomerOrder(ctx)
	})

	req := httptest.NewRequest(http.MethodGet, customerPrefix+"/orders/abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateOrderBadRequest(t *testing.T) {
	mockService := new(mockCustomerService)

	r := gin.New()
	r.POST(customerPrefix+"/orders", func(ctx *gin.Context) {
		ctx.Set("login", "customer1")
		ctx.Set("role", 1)
		handler.NewCustomerHandler(mockService).CreateOrder(ctx)
	})

	req := httptest.NewRequest(http.MethodPost, customerPrefix+"/orders", strings.NewReader(`{"deadline":"not-a-date"}`))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteOrderBadRequest(t *testing.T) {
	mockService := new(mockCustomerService)

	r := gin.New()
	r.DELETE(customerPrefix+"/orders/:id", func(ctx *gin.Context) {
		ctx.Set("login", "customer1")
		ctx.Set("role", 1)
		handler.NewCustomerHandler(mockService).DeleteOrder(ctx)
	})

	req := httptest.NewRequest(http.MethodDelete, customerPrefix+"/orders/abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSaveToBasketBadRequest(t *testing.T) {
	mockService := new(mockCustomerService)

	r := gin.New()
	r.POST(customerPrefix+"/basket", func(ctx *gin.Context) {
		ctx.Set("login", "customer1")
		ctx.Set("role", 1)
		handler.NewCustomerHandler(mockService).SaveToBasket(ctx)
	})

	req := httptest.NewRequest(http.MethodPost, customerPrefix+"/basket", strings.NewReader(`{"itemID":"bad","amount":-1}`))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteFromBasketBadRequest(t *testing.T) {
	mockService := new(mockCustomerService)

	r := gin.New()
	r.DELETE(customerPrefix+"/basket/:id", func(ctx *gin.Context) {
		ctx.Set("login", "customer1")
		ctx.Set("role", 1)
		handler.NewCustomerHandler(mockService).DeleteFromBasket(ctx)
	})

	req := httptest.NewRequest(http.MethodDelete, customerPrefix+"/basket/abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateOrder(t *testing.T) {
	mockService := new(mockCustomerService)
	orderDTO := model.OrderCreationDTO{
		Deadline: time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC),
		Tasks: []model.TaskCreationDTO{
			{ItemID: 1, Amount: 2},
			{ItemID: 2, Amount: 1},
		}}
	mockService.On("CreateOrder", "customer1", 1, orderDTO).Return(1, nil)

	r := gin.New()
	r.POST(customerPrefix+"/orders", func(ctx *gin.Context) {
		ctx.Set("login", "customer1")
		ctx.Set("role", 1)
		handler.NewCustomerHandler(mockService).CreateOrder(ctx)
	})

	body, _ := json.Marshal(orderDTO)
	req := httptest.NewRequest(http.MethodPost, customerPrefix+"/orders", strings.NewReader(string(body)))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "/customer/orders/1", w.Header().Get("Location"))
}

func TestDeleteOrder(t *testing.T) {
	mockService := new(mockCustomerService)
	mockService.On("DeleteOrder", "customer1", 1, 1).Return(nil)

	r := gin.New()
	r.DELETE(customerPrefix+"/orders/:id", func(ctx *gin.Context) {
		ctx.Set("login", "customer1")
		ctx.Set("role", 1)
		handler.NewCustomerHandler(mockService).DeleteOrder(ctx)
	})

	req := httptest.NewRequest(http.MethodDelete, customerPrefix+"/orders/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestGetItems(t *testing.T) {
	mockService := new(mockCustomerService)
	mockItems := []model.Item{
		{Id: 1, Name: "Item 1", Price: 100},
		{Id: 2, Name: "Item 2", Price: 200},
	}
	mockService.On("GetItems", "customer1", 1).Return(mockItems, nil)

	r := gin.New()
	r.GET(customerPrefix+"/items", func(ctx *gin.Context) {
		ctx.Set("login", "customer1")
		ctx.Set("role", 1)
		handler.NewCustomerHandler(mockService).GetItems(ctx)
	})

	req := httptest.NewRequest(http.MethodGet, customerPrefix+"/items", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var responseItems []model.Item
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &responseItems))
	assert.Equal(t, mockItems, responseItems)
}

func TestGetBasket(t *testing.T) {
	mockService := new(mockCustomerService)
	mockBasket := []model.Task{
		{Id: 1, Name: "Task 1", OrderID: 1, ItemID: 1, Amount: 2, Finished: false, Price: 100},
		{Id: 2, Name: "Task 2", OrderID: 1, ItemID: 2, Amount: 1, Finished: false, Price: 50},
	}
	mockService.On("GetBasket", "customer1", 1).Return(mockBasket, nil)

	r := gin.New()
	r.GET(customerPrefix+"/basket", func(ctx *gin.Context) {
		ctx.Set("login", "customer1")
		ctx.Set("role", 1)
		handler.NewCustomerHandler(mockService).GetBasket(ctx)
	})

	req := httptest.NewRequest(http.MethodGet, customerPrefix+"/basket", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var responseBasket []model.Task
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &responseBasket))
	assert.Equal(t, mockBasket, responseBasket)
}

func TestSaveToBasket(t *testing.T) {
	mockService := new(mockCustomerService)
	basketItem := model.TaskCreationDTO{ItemID: 3, Amount: 5}
	mockService.On("SaveToBasket", "customer1", 1, basketItem).Return(nil)

	r := gin.New()
	r.POST(customerPrefix+"/basket", func(ctx *gin.Context) {
		ctx.Set("login", "customer1")
		ctx.Set("role", 1)
		handler.NewCustomerHandler(mockService).SaveToBasket(ctx)
	})

	body, _ := json.Marshal(basketItem)
	req := httptest.NewRequest(http.MethodPost, customerPrefix+"/basket", strings.NewReader(string(body)))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteFromBasket(t *testing.T) {
	mockService := new(mockCustomerService)
	mockService.On("DeleteFromBasket", "customer1", 1, 1).Return(nil)

	r := gin.New()
	r.DELETE(customerPrefix+"/basket/:id", func(ctx *gin.Context) {
		ctx.Set("login", "customer1")
		ctx.Set("role", 1)
		handler.NewCustomerHandler(mockService).DeleteFromBasket(ctx)
	})

	req := httptest.NewRequest(http.MethodDelete, customerPrefix+"/basket/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
