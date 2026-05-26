package main

import (
	"crud-go/internal/handler"
	"crud-go/internal/middleware"
	"crud-go/internal/repository"
	"crud-go/internal/service"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}

	dbConn, connErr := repository.Connect()
	if connErr != nil {
		log.Fatal(connErr)
	}
	defer repository.Close(dbConn)

	basketRepository, basketRepositoryCreationErr := repository.NewBasketRepository(dbConn)
	if basketRepositoryCreationErr != nil {
		log.Fatal(basketRepositoryCreationErr)
	}
	itemRepository, itemRepositoryCreationErr := repository.NewItemRepository(dbConn)
	if itemRepositoryCreationErr != nil {
		log.Fatal(itemRepositoryCreationErr)
	}
	orderRepository, orderRepositoryCreationErr := repository.NewOrderRepository(dbConn)
	if orderRepositoryCreationErr != nil {
		log.Fatal(orderRepositoryCreationErr)
	}
	taskRepository, taskRepositoryCreationErr := repository.NewTaskRepository(dbConn)
	if taskRepositoryCreationErr != nil {
		log.Fatal(taskRepositoryCreationErr)
	}
	tokenRepository, tokenRepositoryCreationErr := repository.NewTokenRepository(dbConn)
	if tokenRepositoryCreationErr != nil {
		log.Fatal(tokenRepositoryCreationErr)
	}
	userRepository, userRepositoryCreationErr := repository.NewUserRepository(dbConn)
	if userRepositoryCreationErr != nil {
		log.Fatal(userRepositoryCreationErr)
	}

	accessLT, convErr := strconv.Atoi(os.Getenv("AUTH_ACCESS_LIFETIME"))
	refreshLT, convErr := strconv.Atoi(os.Getenv("AUTH_REFRESH_LIFETIME"))
	if convErr != nil {
		log.Fatal(convErr)
	}

	authService := service.NewAuthService(userRepository, tokenRepository, accessLT, refreshLT, os.Getenv("AUTH_ACCESS_KEY"), os.Getenv("AUTH_REFRESH_KEY"))
	customerService := service.NewCustomerService(&orderRepository, &taskRepository, basketRepository, itemRepository)
	managerService := service.NewManagerService(&orderRepository, &taskRepository, userRepository)
	workerService := service.NewWorkerService(&orderRepository, &taskRepository)

	authHandler := handler.NewAuthHandler(authService)
	customerHandler := handler.NewCustomerHandler(customerService)
	managerHandler := handler.NewManagerHandler(managerService)
	workerHandler := handler.NewWorkerHandler(workerService)

	authMiddleware := middleware.NewAuthMiddleware(authService)

	router := gin.Default()
	customerGroup := router.Group("/api/v1/customer")
	workerGroup := router.Group("api/v1/worker")
	managerGroup := router.Group("api/v1/manager")
	authGroup := router.Group("api/v1/auth")

	customerGroup.GET("/orders", customerHandler.GetAllOrders)
	customerGroup.GET("/orders/:id", customerHandler.GetCustomerOrder)
	customerGroup.POST("/orders", customerHandler.CreateOrder)
	customerGroup.DELETE("orders/:id", customerHandler.DeleteOrder)
	customerGroup.GET("/items", customerHandler.GetItems)
	customerGroup.GET("/basket", customerHandler.GetBasket)
	customerGroup.POST("/basket", customerHandler.SaveToBasket)
	customerGroup.DELETE("/basket/:id", customerHandler.DeleteFromBasket)
	customerGroup.POST("", authHandler.Register)

	workerGroup.GET("/orders", workerHandler.GetOrdersByWorker)
	workerGroup.GET("/orders/:id", workerHandler.GetOrder)
	workerGroup.POST("/tasks/:id", workerHandler.SetTaskCompleted)

	managerGroup.GET("/orders", managerHandler.GetOrdersByManager)
	managerGroup.POST("/orders/:id", managerHandler.AssignWorkerToOrder)
	managerGroup.GET("/orders/:id", managerHandler.GetOrderByID)
	managerGroup.GET("/workers", managerHandler.GetAllWorkers)

	authGroup.POST("/login", authHandler.Login)
	authGroup.POST("/refresh", authHandler.RefreshToken)
	authGroup.POST("/logout", authHandler.Logout)

	customerGroup.Use(authMiddleware.AuthMiddlewareFunc)
	workerGroup.Use(authMiddleware.AuthMiddlewareFunc)
	managerGroup.Use(authMiddleware.AuthMiddlewareFunc)

	router.Run(":" + os.Getenv("APP_PORT"))
}
