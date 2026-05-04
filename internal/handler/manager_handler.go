package handler

import (
	"crud-go/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ManagerHandler struct {
	managerService service.IManagerService
}

func NewManagerHandler(managerService service.IManagerService) ManagerHandler {
	return ManagerHandler{managerService: managerService}
}

func (mh ManagerHandler) GetOrdersByManager(ctx *gin.Context) {
	//Извлекаем логин и роль менеджера из контекста
	login, role := getUserInfo(ctx)
	if login == "" || role == 0 {
		return
	}
	//Получаем заказы менеджера через сервис по логину и роли
	orders, serviceErr := mh.managerService.GetAllOrders(login, role)
	if serviceErr != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving orders: " + serviceErr.Error()})
		return
	}

	//Возвращаем список заказов в формате JSON
	ctx.JSON(http.StatusOK, orders)
}

func (mh ManagerHandler) GetOrderByID(ctx *gin.Context) {
	//Извлекаем логин и роль менеджера из контекста
	login, role := getUserInfo(ctx)
	if login == "" || role == 0 {
		return
	}
	// Получаем ID заказа из параметров URL и преобразуем его в целое число
	orderID, paramErr := strconv.Atoi(ctx.Param("orderId"))
	if paramErr != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	//Получаем заказ менеджера через сервис по логину, роли и ID заказа
	order, serviceErr := mh.managerService.GetOrderById(login, role, orderID)
	if serviceErr != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving order: " + serviceErr.Error()})
		return
	}
	//Возвращаем заказ в формате JSON
	ctx.JSON(http.StatusOK, order)
}

func (mh ManagerHandler) SetWorkerLogin(ctx *gin.Context) {
	//Извлекаем логин и роль менеджера из контекста
	login, role := getUserInfo(ctx)
	if login == "" || role == 0 {
		return
	}
	// Получаем ID заказа из параметров URL и преобразуем его в целое число
	orderID, paramErr := strconv.Atoi(ctx.Param("orderId"))
	if paramErr != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}
	// Получаем логин рабочего из тела запроса
	var requestBody struct {
		WorkerLogin string `json:"workerLogin"`
	}
	if bindErr := ctx.ShouldBindJSON(&requestBody); bindErr != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + bindErr.Error()})
		return
	}
	//Устанавливаем логин рабочего для заказа через сервис по логину, роли, ID заказа и логину рабочего
	serviceErr := mh.managerService.AssignWorkerToOrder(login, role, orderID, requestBody.WorkerLogin)
	if serviceErr != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error setting worker login: " + serviceErr.Error()})
		return
	}
	//Возвращаем статус 200 OK без тела ответа
	ctx.Status(http.StatusOK)
}
