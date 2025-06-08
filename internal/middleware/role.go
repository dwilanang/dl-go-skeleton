package middleware

import (
	"net/http"

	"slices"

	"github.com/dwilanang/psp/internal/auth/model"
	"github.com/gin-gonic/gin"
)

// RequireRole returns a Gin middleware handler that enforces role-based access control.
//
// The middleware checks if the authenticated user (stored in the Gin context under the key "user")
// has a role included in the allowedRoles list.
//
// If the user is not authenticated or the token claims are invalid, it aborts the request with
// a 401 Unauthorized status.
//
// If the user's role is not in the allowedRoles slice, it aborts with a 403 Forbidden status,
// indicating insufficient privileges.
//
// Parameters:
//   - allowedRoles: a variadic list of strings representing roles permitted to access the endpoint.
//
// Returns:
//   - gin.HandlerFunc: the middleware function that enforces role restrictions.
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		claims, ok := val.(*model.TokenClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		if slices.Contains(allowedRoles, claims.Role) {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden: insufficient privileges"})
	}
}
