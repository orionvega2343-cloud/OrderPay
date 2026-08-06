package postgres

import (
	"OrderPay/pkg/config"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

//Реализация пул соединений,
//с health-check проверкой на валидное соединение

func HealthCheck(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		//Ограничиваем время проверки,
		//чтобы запрос не завис
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		//PingContext проверяет,
		//сохраняется ли соединение с базой данных,
		//и при необходимости устанавливает его
		if err := db.PingContext(ctx); err != nil {
			slog.Error("Error pinging database", "error", err)
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "failed",
				"error":  "postgres connection failed",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	}
}

func ConnDB(cfg config.Config) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s", cfg.DB.Host, cfg.DB.Port, cfg.DB.Username, cfg.DB.Name, cfg.DB.Password, cfg.DB.SslMode)
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}
