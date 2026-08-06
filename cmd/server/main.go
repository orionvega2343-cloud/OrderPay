package main

import (
	"OrderPay/pkg/config"
	"OrderPay/pkg/middleware"
	"OrderPay/pkg/postgres"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	//Вызываем конфиг и бд
	cfg := config.MustLoad()
	conn, err := postgres.ConnDB(*cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		cerr := conn.Close()
		if cerr != nil {
			log.Println(cerr)
		}
	}()

	//Вызов recovery middleware, всегда идет первым,
	//чтобы ловить панику нижних методов/сервисов
	r.Use(middleware.Recovery())
	r.Use(middleware.RequestId())
	r.Use(middleware.Timeout())
	r.Use(middleware.Logger())
	r.GET("/connection", postgres.HealthCheck(conn))

	//TODO: пробросить порт в конфиг
	r.Run(":8080")
}
