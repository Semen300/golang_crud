package handler

import (
	"net/http"
	"strconv"

	"crud-go/internal/service"

	"github.com/gin-gonic/gin"
)

type WorkerHandler struct {
	workerService service.IWorkerService
}

func NewWorkerHandler(workerService service.IWorkerService) WorkerHandler {
	return WorkerHandler{workerService: workerService}
}

// GetOrdersByWorker предназначена для получения всех заказов, назначенных рабочему.
//
// Принимает контекст запроса и сервис для работы с рабочими,
// возвращает список всех заказов, связанных с рабочим, который сделал запрос.
func (wh *WorkerHandler) GetOrdersByWorker(ctx *gin.Context) {
	//Извлекаем логин и роль рабочего из контекста
	login, role := getUserInfo(ctx)
	if login == "" || role == 0 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Incorrect user claims"})
		return
	}
	//Получаем заказы рабочего через сервис по логину и роли
	orders, serviceErr := wh.workerService.GetAllOrders(login, role)
	if serviceErr != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving orders: " + serviceErr.Error()})
		return
	}

	//Возвращаем список заказов в формате JSON
	ctx.JSON(http.StatusOK, orders)
}

// GetOrder предназначена для получения информации о конкретном заказе, назначенном рабочему.
//
// Принимает контекст запроса и сервис для работы с рабочими,
// возвращает информацию о заказе с указанным ID.
func (wh WorkerHandler) GetOrder(ctx *gin.Context) {
	//Извлекаем логин и роль рабочего из контекста
	login, role := getUserInfo(ctx)
	if login == "" || role == 0 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Incorrect user claims"})
		return
	}
	// Получаем ID заказа из параметров URL и преобразуем его в целое число
	orderID, paramErr := strconv.Atoi(ctx.Param("orderId"))
	if paramErr != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	//Получаем заказ рабочего через сервис по логину, роли и ID заказа
	order, serviceErr := wh.workerService.GetOrderById(login, role, orderID)
	if serviceErr != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving order: " + serviceErr.Error()})
		return
	}

	//Возвращаем заказ в формате JSON
	ctx.JSON(http.StatusOK, order)
}

// SetTaskCompleted предназначена для установки статуса задачи как выполненной.
//
// Принимает контекст запроса и сервис для работы с рабочими,
// устанавливает статус задачи как выполненной по ID задачи, указанному в URL, для рабочего, который сделал запрос.
// В случае успеха возвращает статус 204 No Content и URL задачи в заголовке Location, иначе - сообщение об ошибке.
func (wh WorkerHandler) SetTaskCompleted(ctx *gin.Context) {
	//Извлекаем логин и роль рабочего из контекста
	login, role := getUserInfo(ctx)
	if login == "" || role == 0 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Incorrect user claims"})
		return
	}
	// Получаем ID задачи из параметров URL и преобразуем его в целое число
	taskID, paramErr := strconv.Atoi(ctx.Param("taskId"))
	if paramErr != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	//Устанавливаем задачу как выполненную через сервис по логину, роли и ID задачи
	serviceErr := wh.workerService.SetTaskCompleted(login, role, taskID)
	if serviceErr != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error setting task as completed: " + serviceErr.Error()})
		return
	}

	//Устанавливаем параметр Location в заголовке ответа, указывая на URL задачи, которая была установлена как выполненная
	ctx.Header("Location", "/tasks/"+strconv.Itoa(taskID))
	ctx.Status(http.StatusNoContent)
}
