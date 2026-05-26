package service_test

/*
import (
	"crud-go/internal/model"
	"crud-go/internal/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewAuthService(t *testing.T) {
	userRepoMock := new(userRepoMock)
	tokenRepoMock := new(tokenRepoMock)
	authService := service.NewAuthService(userRepoMock, tokenRepoMock, 1000, 10000, "accessSecret", "refreshSecret")
	assert.NotEqual(t, service.AuthService{}, authService)
}

func TestRegisterNewCustomer(t *testing.T) {
	userRepoMock := new(userRepoMock)
	testService := service.NewAuthService(userRepoMock, nil, 1000, 10000, "accessSecret", "refreshSecret")

	customerToSave := model.NewCustomer("customer1", "1111", "AAAA", "89111111111", "aaaa@gmail.com")
	userRepoMock.On("SaveCustomer", customerToSave).Return(nil)

	saveErr := testService.RegisterNewCustomer(
		customerToSave.Login,
		customerToSave.Password,
		customerToSave.Fio,
		customerToSave.Number,
		customerToSave.Email,
	)

	assert.NoError(t, saveErr)
}

func TestLogin(t *testing.T) {
	userRepoMock := new(userRepoMock)
	tokenRepoMock := new(tokenRepoMock)
	testService := service.NewAuthService(userRepoMock, tokenRepoMock, 1000, 10000, "accessSecret", "refreshSecret")

	testLogin := "customer1"
	testPassword := "password1"
	hashedPassword := service.CreateHash(testPassword)
	tokenToSave := model.RefreshToken{
		TokenID:   "1111",
		Login:     testLogin,
		TokenHash: "11111",
		ExpiresAt: time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC),
	}

	userRepoMock.On("GetRoleByLogin", testLogin).Return(1, hashedPassword, nil)
	tokenRepoMock.On("Save", tokenToSave.TokenID, tokenToSave.Login, tokenToSave.TokenHash, tokenToSave.ExpiresAt).Return(nil)
	tokenRepoMock.On("Revoke", testLogin).Return(nil)
	claims := model.Claims{Login: testLogin, Role: 1}

	accessClaims, accessToken, refreshToken, loginErr := testService.Login(testLogin, testPassword)
	expectedAT, _ := testService.GenerateAccessToken(claims)
	_, expectedRT, _, _ := testService.GenerateRefreshToken(claims)

	assert.NoError(t, loginErr)
	assert.Equal(t, claims, accessClaims)
	assert.Equal(t, expectedAT, accessToken)
	assert.Equal(t, expectedRT, refreshToken)
}
*/
