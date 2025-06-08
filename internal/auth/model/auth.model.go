package model

import "github.com/golang-jwt/jwt/v5"

type TokenClaims struct {
	UUID string `json:"sub"`
	ID   int64  `json:"uid"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}
