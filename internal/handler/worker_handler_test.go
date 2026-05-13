package handler_test

import (
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

const workerPrefix = "/api/v1/worker"

type mockWorkerService struct {
	mock.Mock
}

func (m *mockWorkerService) GetAllOrders(login string, role int) ([]model.Order, error) {
	args := m.Called(login, role)
	return args.Get(0).([]model.Order), args.Error(1)
}

func (m *mockWorkerService) GetOrderById(login string, role int, id int) (model.Order, error) {
	args := m.Called(login, role, id)
	return args.Get(0).(model.Order), args.Error(1)
}

func (m *mockWorkerService) SetTaskCompleted(login string, role int, id int) error {
	args := m.Called(login, role, id)
	return args.Error(0)
}

func TestGetOrdersByWorker200(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mockWorkerService)
	mockOrders := []model.Order{
		{ID: 1, Name: "Order 1", Deadline: time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC), ManagerLogin: "manager1", WorkerLogin: "worker1", CustomerLogin: "customer1", PercentOfComplition: 50.0, PriceTotal: 100, PriceUnfinished: 50, Status: 2},
		{ID: 2, Name: "Order 2", Deadline: time.Date(2026, time.May, 1, 0, 0, 0, 0, time.UTC), ManagerLogin: "manager2", WorkerLogin: "worker1", CustomerLogin: "customer1", PercentOfComplition: 100.0, PriceTotal: 200, PriceUnfinished: 0, Status: 3},
	}
	mockService.On("GetAllOrders", "worker1", 2).Return(mockOrders, nil)

	req := httptest.NewRequest(http.MethodGet, workerPrefix+"/orders", nil)
	w := httptest.NewRecorder()
	testCtx, _ := gin.CreateTestContext(w)
	testCtx.Request = req
	testCtx.Set("login", "worker1")
	testCtx.Set("role", 2)

	testHandler := handler.NewWorkerHandler(mockService)
	testHandler.GetOrdersByWorker(testCtx)

	assert.Equal(t, http.StatusOK, w.Code)
	var responseOrders []model.Order
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &responseOrders))
	assert.Equal(t, mockOrders, responseOrders)
}

func TestGetOrder200(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mockWorkerService)
	mockOrder := model.Order{ID: 1, Name: "Order 1", Deadline: time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC), ManagerLogin: "manager1", WorkerLogin: "worker1", CustomerLogin: "customer1", PercentOfComplition: 50.0, PriceTotal: 100, PriceUnfinished: 50, Status: 2}
	mockService.On("GetOrderById", "worker1", 2, 1).Return(mockOrder, nil)

	r := gin.New()
	r.GET(workerPrefix+"/orders/:orderId", func(ctx *gin.Context) {
		ctx.Set("login", "worker1")
		ctx.Set("role", 2)
		handler.NewWorkerHandler(mockService).GetOrder(ctx)
	})

	req := httptest.NewRequest(http.MethodGet, workerPrefix+"/orders/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var responseOrder model.Order
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &responseOrder))
	assert.Equal(t, mockOrder, responseOrder)
}

func TestGetOrderBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mockWorkerService)

	r := gin.New()
	r.GET(workerPrefix+"/orders/:orderId", func(ctx *gin.Context) {
		ctx.Set("login", "worker1")
		ctx.Set("role", 2)
		handler.NewWorkerHandler(mockService).GetOrder(ctx)
	})

	req := httptest.NewRequest(http.MethodGet, workerPrefix+"/orders/abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSetTaskCompleted204(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mockWorkerService)
	mockService.On("SetTaskCompleted", "worker1", 2, 1).Return(nil)

	r := gin.New()
	r.PUT(workerPrefix+"/tasks/:taskId", func(ctx *gin.Context) {
		ctx.Set("login", "worker1")
		ctx.Set("role", 2)
		handler.NewWorkerHandler(mockService).SetTaskCompleted(ctx)
	})

	req := httptest.NewRequest(http.MethodPut, workerPrefix+"/tasks/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "/tasks/1", w.Header().Get("Location"))
}

func TestSetTaskCompletedBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mockWorkerService)

	r := gin.New()
	r.PUT(workerPrefix+"/tasks/:taskId", func(ctx *gin.Context) {
		ctx.Set("login", "worker1")
		ctx.Set("role", 2)
		handler.NewWorkerHandler(mockService).SetTaskCompleted(ctx)
	})

	req := httptest.NewRequest(http.MethodPut, workerPrefix+"/tasks/abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
