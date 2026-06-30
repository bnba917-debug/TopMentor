package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	jwtpkg "github.com/topmentor/backend/pkg/jwt"
	"github.com/topmentor/backend/pkg/response"
)

const ContextUserIDKey = "user_id"

func Auth(jwtMgr *jwtpkg.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			response.Fail(c, http.StatusUnauthorized, 40101, "未登录或 Token 失效")
			c.Abort()
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		claims, err := jwtMgr.Parse(token)
		if err != nil {
			response.Fail(c, http.StatusUnauthorized, 40101, "未登录或 Token 失效")
			c.Abort()
			return
		}

		c.Set(ContextUserIDKey, claims.UserID)
		if claims.MentorID > 0 {
			c.Set(ContextMentorIDKey, claims.MentorID)
		}
		c.Next()
	}
}

func GetUserID(c *gin.Context) (int64, bool) {
	v, ok := c.Get(ContextUserIDKey)
	if !ok {
		return 0, false
	}
	userID, ok := v.(int64)
	return userID, ok
}

func GetMentorID(c *gin.Context) (int64, bool) {
	v, ok := c.Get(ContextMentorIDKey)
	if !ok {
		return 0, false
	}
	mentorID, ok := v.(int64)
	return mentorID, ok
}
