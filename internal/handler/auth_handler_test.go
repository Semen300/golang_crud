package handler_test

import (
	"bytes"
	"crud-go/internal/handler"
	"crud-go/internal/model"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const authPrefix = "/api/v1/auth"

type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) RegisterNewCustomer(login, password, fio, number, email string) error {
	args := m.Called(login, password, fio, number, email)
	return args.Error(0)
}

func (m *mockAuthService) Login(login, password string) (int, string, string, error) {
	args := m.Called(login, password)
	return args.Int(0), args.String(1), args.String(2), args.Error(3)
}

func (m *mockAuthService) Refresh(refreshToken string) (string, error) {
	args := m.Called(refreshToken)
	return args.String(0), args.Error(1)
}

func (m *mockAuthService) Logout(refreshToken string) error {
	args := m.Called(refreshToken)
	return args.Error(0)
}

func (m *mockAuthService) GenerateAccessToken(claims model.Claims) (string, error) {
	args := m.Called(claims)
	return args.String(0), args.Error(1)
}

func (m *mockAuthService) GenerateRefreshToken(claims model.Claims) (string, string, time.Time, error) {
	args := m.Called(claims)
	return args.String(0), args.String(1), args.Get(2).(time.Time), args.Error(3)
}

func (m *mockAuthService) ParseAccessToken(accessToken string) (model.Claims, error) {
	args := m.Called(accessToken)
	return args.Get(0).(model.Claims), args.Error(1)
}

func (m *mockAuthService) ParseRefreshToken(refreshToken string) (model.Claims, error) {
	args := m.Called(refreshToken)
	return args.Get(0).(model.Claims), args.Error(1)
}

func TestRegister200(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mockAuthService)
	mockService.On("RegisterNewCustomer", "testlogin", "testpass", "Test Fio", "123456", "test@email.com").Return(nil)

	registrationData := model.NewCustomer(
		"testlogin",
		"testpass",
		"Test Fio",
		"123456",
		"test@email.com",
	)
	body, _ := json.Marshal(registrationData)
	req := httptest.NewRequest(http.MethodPost, authPrefix+"/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testCtx, _ := gin.CreateTestContext(w)
	testCtx.Request = req

	testHandler := handler.NewAuthHandler(mockService)
	testHandler.Register(testCtx)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, "Registration successful", response["message"])
	mockService.AssertExpectations(t)
}

func TestRegisterBadRequestInvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mockAuthService)

	req := httptest.NewRequest(http.MethodPost, authPrefix+"/register", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testCtx, _ := gin.CreateTestContext(w)
	testCtx.Request = req

	testHandler := handler.NewAuthHandler(mockService)
	testHandler.Register(testCtx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]string
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, "Invalid request body", response["error"])
}

func TestRegister500(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mockAuthService)
	mockService.On("RegisterNewCustomer", "testlogin", "testpass", "Test Fio", "123456", "test@email.com").Return(errors.New("registration failed"))

	registrationData := model.NewCustomer(
		"testlogin",
		"testpass",
		"Test Fio",
		"123456",
		"test@email.com",
	)
	body, _ := json.Marshal(registrationData)
	req := httptest.NewRequest(http.MethodPost, authPrefix+"/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testCtx, _ := gin.CreateTestContext(w)
	testCtx.Request = req

	testHandler := handler.NewAuthHandler(mockService)
	testHandler.Register(testCtx)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var response map[string]string
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Contains(t, response["error"], "Error during registration")
	mockService.AssertExpectations(t)
}

func TestLogin200(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mockAuthService)
	mockService.On("Login", "testlogin", "testpass").Return(1, "access_token", "refresh_token", nil)

	credentials := map[string]string{
		"login":    "testlogin",
		"password": "testpass",
	}
	body, _ := json.Marshal(credentials)
	req := httptest.NewRequest(http.MethodPost, authPrefix+"/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testCtx, _ := gin.CreateTestContext(w)
	testCtx.Request = req

	testHandler := handler.NewAuthHandler(mockService)
	testHandler.Login(testCtx)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, float64(1), response["role"])
	assert.Equal(t, "access_token", response["accessToken"])
	assert.Equal(t, "refresh_token", response["refreshToken"])
	assert.Equal(t, float64(900), response["expiresIn"])
	mockService.AssertExpectations(t)
}

func TestLoginBadRequestInvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mockAuthService)

	req := httptest.NewRequest(http.MethodPost, authPrefix+"/login", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testCtx, _ := gin.CreateTestContext(w)
	testCtx.Request = req

	testHandler := handler.NewAuthHandler(mockService)
	testHandler.Login(testCtx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]string
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, "Invalid request body", response["error"])
}

func TestLogin401(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mockAuthService)
	mockService.On("Login", "testlogin", "testpass").Return(0, "", "", errors.New("invalid credentials"))

	credentials := map[string]string{
		"login":    "testlogin",
		"password": "testpass",
	}
	body, _ := json.Marshal(credentials)
	req := httptest.NewRequest(http.MethodPost, authPrefix+"/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testCtx, _ := gin.CreateTestContext(w)
	testCtx.Request = req

	testHandler := handler.NewAuthHandler(mockService)
	testHandler.Login(testCtx)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var response map[string]string
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Contains(t, response["error"], "Invalid login or password")
	mockService.AssertExpectations(t)
}

func TestRefreshToken200(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mockAuthService)
	mockService.On("RefreshToken", "refresh_token").Return("new_access_token", nil)

	request := map[string]string{
		"refreshToken": "refresh_token",
	}
	body, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, authPrefix+"/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testCtx, _ := gin.CreateTestContext(w)
	testCtx.Request = req

	testHandler := handler.NewAuthHandler(mockService)
	testHandler.RefreshToken(testCtx)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, "new_access_token", response["accessToken"])
	assert.Equal(t, float64(900), response["expiresIn"])
	mockService.AssertExpectations(t)
}

func TestRefreshTokenBadRequestInvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mockAuthService)

	req := httptest.NewRequest(http.MethodPost, authPrefix+"/refresh", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testCtx, _ := gin.CreateTestContext(w)
	testCtx.Request = req

	testHandler := handler.NewAuthHandler(mockService)
	testHandler.RefreshToken(testCtx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]string
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, "Invalid request body", response["error"])
}

func TestRefreshToken401(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mockAuthService)
	mockService.On("RefreshToken", "refresh_token").Return("", errors.New("invalid token"))

	request := map[string]string{
		"refreshToken": "refresh_token",
	}
	body, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, authPrefix+"/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testCtx, _ := gin.CreateTestContext(w)
	testCtx.Request = req

	testHandler := handler.NewAuthHandler(mockService)
	testHandler.RefreshToken(testCtx)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var response map[string]string
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Contains(t, response["error"], "Invalid refresh token")
	mockService.AssertExpectations(t)
}

func TestLogout200(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mockAuthService)
	mockService.On("Logout", "refresh_token").Return(nil)

	request := map[string]string{
		"refreshToken": "refresh_token",
	}
	body, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, authPrefix+"/logout", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testCtx, _ := gin.CreateTestContext(w)
	testCtx.Request = req

	testHandler := handler.NewAuthHandler(mockService)
	testHandler.Logout(testCtx)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, "Successfully logged out", response["message"])
	mockService.AssertExpectations(t)
}

func TestLogoutBadRequestInvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mockAuthService)

	req := httptest.NewRequest(http.MethodPost, authPrefix+"/logout", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testCtx, _ := gin.CreateTestContext(w)
	testCtx.Request = req

	testHandler := handler.NewAuthHandler(mockService)
	testHandler.Logout(testCtx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]string
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, "Invalid request body", response["error"])
}

func TestLogout500(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mockAuthService)
	mockService.On("Logout", "refresh_token").Return(errors.New("logout failed"))

	request := map[string]string{
		"refreshToken": "refresh_token",
	}
	body, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, authPrefix+"/logout", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testCtx, _ := gin.CreateTestContext(w)
	testCtx.Request = req

	testHandler := handler.NewAuthHandler(mockService)
	testHandler.Logout(testCtx)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var response map[string]string
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Contains(t, response["error"], "Error during logout")
	mockService.AssertExpectations(t)
}
