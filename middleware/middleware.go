package middleware

import (
	"KnowEase/services"
	"net/http"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Middleware struct {
	TokenService *services.TokenService
}

func NewMiddleWare(TokenService *services.TokenService) *Middleware {
	return &Middleware{TokenService: TokenService}
}
func (m *Middleware) Verifytoken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求头中的 token
		tokenString := c.GetHeader("Authorization")
		role, _ := m.TokenService.VerifyToken(tokenString)
		if role != "user" {
			c.JSON(http.StatusForbidden, gin.H{"message": "身份检验失败！"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// 解决跨域问题
func (m *Middleware) Cors() gin.HandlerFunc {
	c := cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:    []string{"Content-Type", "Access-Token", "Authorization"},
		MaxAge:          6 * time.Hour,
	}

	return cors.New(c)
}

