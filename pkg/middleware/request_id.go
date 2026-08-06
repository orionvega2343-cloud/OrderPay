package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ctxKey string

const requestIDKey ctxKey = "request_id"

//Генерирует/пробрасывает сквозной id запроса,
//если входящий X-Request-ID уже есть - переиспользуем, если нет, создаем новый с помощью uuid, кладем id:
//в context для пробрасывания вниз по цепочке вызовов и для обработчика внутри сервиса,
//в заголовок ответа чтобы клиент/поддержка видели id в ответе

func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Request.Header.Get("X-Request-Id")
		if id == "" {
			id = uuid.NewString()
		}

		ctx := context.WithValue(c.Request.Context(), requestIDKey, id)

		c.Request = c.Request.WithContext(ctx)
		c.Header("X-Request-Id", id)
		c.Next()
	}
}
