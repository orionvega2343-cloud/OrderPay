package main

import (
	orderHandler "OrderPay/internal/order/handler"
	orderRepos "OrderPay/internal/order/repository"
	orderService "OrderPay/internal/order/service"
	paymentHandler "OrderPay/internal/payment/handler"
	paymentRepos "OrderPay/internal/payment/repository"
	paymentService "OrderPay/internal/payment/service"
	"OrderPay/pkg/config"
	"OrderPay/pkg/logger"
	"OrderPay/pkg/middleware"
	"OrderPay/pkg/postgres"
	"OrderPay/pkg/transaction"
	"log"
	"log/slog"

	_ "OrderPay/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title OrderPay API
// @version 1.0
// @description Сервис управления заказами, платежами и возвратами.
// @BasePath /api

func main() {
	//Структурированный JSON лог
	slog.SetDefault(logger.NewLogger())
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

	transactor := transaction.NewTransactor(conn)
	//Репозитории
	orderRepo := orderRepos.NewOrderRepo(conn)
	paymentRepo := paymentRepos.NewPaymentRepo(conn)
	refundRepo := paymentRepos.NewRefundRepo(conn)

	//Сервисы
	orderSvc := orderService.NewOrderService(orderRepo, transactor)
	paymentSvc := paymentService.NewPaymentService(paymentRepo, transactor, orderRepo)
	refundSvc := paymentService.NewRefundService(paymentRepo, refundRepo)
	aggregationSvc := paymentService.NewAggregationService(paymentRepo)

	//HTTP-обработчики
	orderHndlr := orderHandler.NewOrderHandlerImpl(orderSvc)
	paymentHndlr := paymentHandler.NewPaymentHandlerImpl(paymentSvc)
	refundHndlr := paymentHandler.NewRefundHandler(refundSvc)
	aggregationHndlr := paymentHandler.NewAggregationHandlerImpl(aggregationSvc)

	//Вызов recovery middleware, всегда идет первым,
	//чтобы ловить панику нижних методов/сервисов
	r.Use(middleware.Recovery())
	r.Use(middleware.RequestId())
	r.Use(middleware.Timeout(cfg.Server.Timeout))
	r.Use(middleware.Logger())
	r.GET("/connection", postgres.HealthCheck(conn))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//Маршруты
	api := r.Group("/api")

	orders := api.Group("/orders")
	{
		orders.POST("", orderHndlr.PostOrder)
		orders.GET("", orderHndlr.GetOrders)
		orders.GET("/:id", orderHndlr.GetOrderByID)
		orders.PATCH("/:id", orderHndlr.UpdateOrder)
		orders.DELETE("/:id", orderHndlr.DeleteOrder)

		// :id здесь - id заказа, платеж создается по заказу
		orders.POST("/:id/payments", paymentHndlr.CreatePayment)
	}

	payments := api.Group("/payments")
	{
		payments.GET("", paymentHndlr.GetAllPayments)
		payments.GET("/total-amount", aggregationHndlr.GetPaymentTotalAmount)
		payments.GET("/:id", paymentHndlr.GetPaymentById)
		payments.PUT("/:id", paymentHndlr.UpdatePayment)
		payments.DELETE("/:id", paymentHndlr.DeletePayment)

		// :id здесь - id платежа, возврат создается по платежу
		payments.POST("/:id/refunds", refundHndlr.CreateRefund)
	}

	refunds := api.Group("/refunds")
	{
		refunds.GET("", refundHndlr.GetAllRefund)
	}

	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal(err)
	}
}
