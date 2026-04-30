package handler

import (
	"crud-go/internal/model"
	"crud-go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.IAuthService
}

func NewAuthHandler(authService service.IAuthService) AuthHandler {
	return AuthHandler{authService: authService}
}

// Register служит для регистрации нового пользователя.
//
// Принимает контекст запроса и сервис для работы с аутентификацией,
// извлекает данные регистрации из тела запроса и передаёт их в сервис для создания нового клиента.
// В случае успеха возвращает сообщение об успешной регистрации, иначе - сообщение об ошибке.
func (ah *AuthHandler) Register(ctx *gin.Context) {
	var registrationData model.Customer
	if err := ctx.ShouldBindJSON(&registrationData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	serviceErr := ah.authService.RegisterNewCustomer(registrationData.Login, registrationData.Password, registrationData.Fio, registrationData.Number, registrationData.Email)
	if serviceErr != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error during registration: " + serviceErr.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

// Login служит для аутентификации пользователя и выдачи ему токенов доступа.
//
// Принимает контекст запроса и сервис для работы с аутентификацией,
// извлекает данные аутентификации из тела запроса и передаёт их в сервис для проверки.
// В случае успеха возвращает роль пользователя, токен доступа, токен обновления и время жизни токена доступа, иначе - сообщение об ошибке.
func (ah *AuthHandler) Login(ctx *gin.Context) {
	var credentials struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	role, accessToken, RefreshToken, serviceErr := ah.authService.Login(credentials.Login, credentials.Password)
	if serviceErr != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid login or password: " + serviceErr.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"role":         role,
		"accessToken":  accessToken,
		"refreshToken": RefreshToken,
		"expiresIn":    900})
}

// RefreshToken служит для обновления токена доступа с помощью токена обновления.
//
// Принимает контекст запроса и сервис для работы с аутентификацией,
// извлекает токен обновления из тела запроса и передаёт его в сервис для проверки.
// В случае успеха возвращает новый токен доступа и время его жизни, иначе - сообщение об ошибке.
func (ah *AuthHandler) RefreshToken(ctx *gin.Context) {
	var request struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	accessToken, serviceErr := ah.authService.RefreshToken(request.RefreshToken)
	if serviceErr != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token: " + serviceErr.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"accessToken": accessToken,
		"expiresIn":   900})
}

// Logout служит для выхода пользователя из системы.
//
// Принимает контекст запроса и сервис для работы с аутентификацией,
// извлекает токен обновления из тела запроса и передаёт его в сервис для проверки.
// В случае успеха возвращает сообщение об успешном выходе, иначе - сообщение об ошибке.
func (ah *AuthHandler) Logout(ctx *gin.Context) {
	var request struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	serviceErr := ah.authService.Logout(request.RefreshToken)
	if serviceErr != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error during logout: " + serviceErr.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}
