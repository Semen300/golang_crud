package handler_test

import (
	"bytes"
	"crud-go/internal/handler"
	"crud-go/internal/model"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const managerPrefix = "/api/v1/manager"

type mockManagerService struct {
	mock.Mock
}

func (m *mockManagerService) GetAllOrders(login string, role int) ([]model.Order, error) {
	args := m.Called(login, role)
	return args.Get(0).([]model.Order), args.Error(1)
}

func (m *mockManagerService) GetOrderById(login string, role int, id int) (model.Order, error) {
	args := m.Called(login, role, id)
	return args.Get(0).(model.Order), args.Error(1)
}

func (m *mockManagerService) AssignWorkerToOrder(login string, role int, orderID int, workerLogin string) error {
	args := m.Called(login, role, orderID, workerLogin)
	return args.Error(0)
}

func TestGetOrdersByManager200(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mockManagerService)
	mockOrders := []model.Order{
		{ID: 1, Name: "Order 1", Deadline: time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC), ManagerLogin: "manager1", WorkerLogin: "worker1", CustomerLogin: "customer1", PercentOfComplition: 50.0, PriseTotal: 100, PriceUnfinished: 50, Status: 2},
		{ID: 2, Name: "Order 2", Deadline: time.Date(2026, time.May, 1, 0, 0, 0, 0, time.UTC), ManagerLogin: "manager1", WorkerLogin: "worker2", CustomerLogin: "customer2", PercentOfComplition: 75.0, PriseTotal: 200, PriceUnfinished: 50, Status: 2},
	}
	mockService.On("GetAllOrders", "manager1", 3).Return(mockOrders, nil)

	req := httptest.NewRequest(http.MethodGet, managerPrefix+"/orders", nil)
	w := httptest.NewRecorder()
	testCtx, _ := gin.CreateTestContext(w)
	testCtx.Request = req
	testCtx.Set("login", "manager1")
	testCtx.Set("role", 3)

	testHandler := handler.NewManagerHandler(mockService)
	testHandler.GetOrdersByManager(testCtx)

	assert.Equal(t, http.StatusOK, w.Code)
	var responseOrders []model.Order
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &responseOrders))
	assert.Equal(t, mockOrders, responseOrders)
}

func TestGetOrderByManager200(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mockManagerService)
	mockOrder := model.Order{ID: 1, Name: "Order 1", Deadline: time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC), ManagerLogin: "manager1", WorkerLogin: "worker1", CustomerLogin: "customer1", PercentOfComplition: 50.0, PriseTotal: 100, PriceUnfinished: 50, Status: 2}
	mockService.On("GetOrderByID", "manager1", 3, 1).Return(mockOrder, nil)

	r := gin.New()
	r.GET(managerPrefix+"/orders/:orderId", func(ctx *gin.Context) {
		ctx.Set("login", "manager1")
		ctx.Set("role", 3)
		handler.NewManagerHandler(mockService).GetOrderByID(ctx)
	})

	req := httptest.NewRequest(http.MethodGet, managerPrefix+"/orders/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var responseOrder model.Order
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &responseOrder))
	assert.Equal(t, mockOrder, responseOrder)
}

func TestGetOrderByManagerBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mockManagerService)

	r := gin.New()
	r.GET(managerPrefix+"/orders/:orderId", func(ctx *gin.Context) {
		ctx.Set("login", "manager1")
		ctx.Set("role", 3)
		handler.NewManagerHandler(mockService).GetOrderByID(ctx)
	})

	req := httptest.NewRequest(http.MethodGet, managerPrefix+"/orders/abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSetWorkerLogin200(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mockManagerService)
	mockService.On("SetWorkerLogin", "manager1", 3, 1, "newworker").Return(nil)

	r := gin.New()
	r.PUT(managerPrefix+"/orders/:orderId/worker", func(ctx *gin.Context) {
		ctx.Set("login", "manager1")
		ctx.Set("role", 3)
		handler.NewManagerHandler(mockService).SetWorkerLogin(ctx)
	})

	requestBody := map[string]string{"workerLogin": "newworker"}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPut, managerPrefix+"/orders/1/worker", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestSetWorkerLoginBadRequestInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mockManagerService)

	r := gin.New()
	r.PUT(managerPrefix+"/orders/:orderId/worker", func(ctx *gin.Context) {
		ctx.Set("login", "manager1")
		ctx.Set("role", 3)
		handler.NewManagerHandler(mockService).SetWorkerLogin(ctx)
	})

	requestBody := map[string]string{"workerLogin": "newworker"}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPut, managerPrefix+"/orders/abc/worker", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSetWorkerLoginBadRequestInvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mockManagerService)

	r := gin.New()
	r.PUT(managerPrefix+"/orders/:orderId/worker", func(ctx *gin.Context) {
		ctx.Set("login", "manager1")
		ctx.Set("role", 1)
		handler.NewManagerHandler(mockService).SetWorkerLogin(ctx)
	})

	req := httptest.NewRequest(http.MethodPut, managerPrefix+"/orders/1/worker", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
