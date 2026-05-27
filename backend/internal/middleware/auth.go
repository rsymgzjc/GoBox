package middleware

import (
	"strings"

	"gobox/backend/internal/config"
	jwtpkg "gobox/backend/pkg/jwt"
	"gobox/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

const userContextKey = "currentUser"

type UserContext struct {
	ID    uint
	Email string
	Role  string
}

func Auth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			response.Fail(c, 401, "未授权", "missing bearer token")
			c.Abort()
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		claims, err := jwtpkg.Parse(cfg.Auth.JWTSecret, token)
		if err != nil {
			response.Fail(c, 401, "令牌无效", err.Error())
			c.Abort()
			return
		}

		c.Set(userContextKey, UserContext{ID: claims.UserID, Email: claims.Email, Role: claims.Role})
		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := CurrentUser(c)
		if !ok || user.Role != "admin" {
			response.Fail(c, 403, "权限不足", "admin only")
			c.Abort()
			return
		}
		c.Next()
	}
}

func CurrentUser(c *gin.Context) (UserContext, bool) {
	raw, ok := c.Get(userContextKey)
	if !ok {
		return UserContext{}, false
	}
	user, ok := raw.(UserContext)
	return user, ok
}
