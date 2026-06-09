package middleware

import (
	"crud-go/internal/service"
	"net/http"
	"strconv"
	"strings"

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
	authHeader := ctx.GetHeader("Authorization")
	token, hasPrefix := strings.CutPrefix(authHeader, "Bearer ")
	if !hasPrefix {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "'Bearer' prefix required at 'Authorization' feild"})
	}
	if authHeader == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: access token required"})
		return
	}

	claims, accessTokenErr := as.ParseAccessToken(token)

	if accessTokenErr != nil {
		if accessTokenErr == jwt.ErrTokenExpired {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: access token expired"})
			return
		} else {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": accessTokenErr.Error()})
			return
		}
	}

	ctx.AddParam("login", claims.Login)
	ctx.AddParam("role", strconv.Itoa(claims.Role))
	ctx.Next()
}
