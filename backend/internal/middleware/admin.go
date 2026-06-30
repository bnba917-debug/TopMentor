package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	jwtpkg "github.com/topmentor/backend/pkg/jwt"
	"github.com/topmentor/backend/pkg/response"
)

func AdminAuth(jwtMgr *jwtpkg.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			response.Fail(c, http.StatusUnauthorized, 40101, "未登录或 Token 失效")
			c.Abort()
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		claims, err := jwtMgr.Parse(token)
		if err != nil || claims.Role != jwtpkg.RoleAdmin {
			response.Fail(c, http.StatusUnauthorized, 40101, "未登录或 Token 失效")
			c.Abort()
			return
		}

		c.Set("admin_username", claims.Subject)
		c.Next()
	}
}
