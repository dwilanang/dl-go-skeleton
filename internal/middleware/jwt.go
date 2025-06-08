package middleware

import (
	"net/http"

	"github.com/dwilanang/psp/internal/auth/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuthMiddleware returns a Gin middleware handler that performs JWT authentication.
//
// It expects the incoming HTTP request to have an "Authorization" header with the format:
//
//	"Bearer <token>"
//
// The middleware extracts the JWT token from the header, parses it using the provided secret key,
// and validates the token's signature and claims.
//
// If the token is missing, malformed, or invalid, the middleware aborts the request with a 401 Unauthorized status
// and a JSON error message.
//
// On successful validation, the parsed token claims are saved into the Gin context with the key "user",
// allowing subsequent handlers to access authenticated user information.
//
// Parameters:
//   - secret: the secret key string used to validate the JWT token signature.
//
// Returns:
//   - gin.HandlerFunc: the middleware function to be used in Gin routes.
func JWTAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid token"})
			return
		}

		tokenStr := authHeader[7:]
		claims := &model.TokenClaims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// save struct on context
		c.Set("user", claims)
		c.Next()
	}
}
