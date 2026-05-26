package middleware

import (
	"crud-go/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	AuthService *service.IAuthService
}

func NewAuthMiddleware(as service.IAuthService) AuthMiddleware {
	return AuthMiddleware{AuthService: &as}
}

func (am AuthMiddleware) AuthMiddlewareFunc(ctx *gin.Context) {
	as := *am.AuthService
	token := ctx.GetHeader("Authorisation")
	if token == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: access token required"})
	}

	claims, accessTokenErr := as.ParseAccessToken(token)

	if accessTokenErr != nil {
		if accessTokenErr == jwt.ErrTokenExpired {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: access token expired"})
		} else {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": accessTokenErr.Error()})
		}
	}

	ctx.AddParam("login", claims.Login)
	ctx.AddParam("role", strconv.Itoa(claims.Role))
	ctx.Next()
}
