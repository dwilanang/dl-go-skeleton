package util

import (
	"errors"

	"github.com/dwilanang/psp/internal/auth/model"
	"github.com/gin-gonic/gin"
)

func GetClaimsID(c *gin.Context) (int64, error) {
	id := int64(0)
	val, exists := c.Get("user")
	if exists {
		claims, ok := val.(*model.TokenClaims)
		if !ok {
			return id, errors.New("error: Invalid token claims")
		}
		id = claims.ID
	}

	return id, nil
}
