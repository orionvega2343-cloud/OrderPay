package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

//Logger пишет структурированный лог по каждому запросу: method,
//path, status, duration_ms, использует stdlib log/slog - без
//zap/zerolog, статус по умолчанию http.StatusOK, если handler не вызвал WriteHeader, go все равно вернет 200
//и мы должны это отразить в логе

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		//Передаем управление обработчику
		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		if query != "" {
			path = path + "?" + query
		}

		slog.Info("http",
			"status", status,
			"latency", latency,
			"clientIP", clientIP,
			"method", method,
			"path", path,
			"query", query,
		)
	}
}
