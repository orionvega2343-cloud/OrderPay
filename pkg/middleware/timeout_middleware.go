package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

func Timeout() gin.HandlerFunc {
	return func(c *gin.Context) {
		//TODO: пробросить динамечское время из конфига
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
