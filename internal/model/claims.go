package model

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	Login string `json:"login"`
	Role  int    `json:"role"`

	jwt.RegisteredClaims
}
