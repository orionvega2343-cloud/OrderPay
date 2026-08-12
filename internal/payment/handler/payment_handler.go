package handler

import (
	"OrderPay/internal/payment/domain"
	"OrderPay/internal/payment/dto"
	"OrderPay/internal/payment/service"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// domain используется только как тип ответа в swagger-аннотациях
var _ domain.Payment

type PaymentHandler interface {
	CreatePayment(c *gin.Context)
	GetPaymentById(c *gin.Context)
	GetAllPayments(c *gin.Context)
	UpdatePayment(c *gin.Context)
	DeletePayment(c *gin.Context)
}

type PaymentHandlerImpl struct {
	svc service.PaymentService
}

func NewPaymentHandlerImpl(svc service.PaymentService) *PaymentHandlerImpl {
	return &PaymentHandlerImpl{svc: svc}
}

// CreatePayment создает платеж по заказу
// @Summary Создать платеж
// @Description id в пути - это id заказа, платеж создается по существующему заказу в статусе created
// @Tags payments
// @Accept json
// @Produce json
// @Param id path int true "ID заказа"
// @Param request body dto.PaymentRequest true "Данные платежа"
// @Success 200 {object} map[string]dto.PaymentResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders/{id}/payments [post]
func (p *PaymentHandlerImpl) CreatePayment(c *gin.Context) {
	id := c.Param("id")
	parsedId, err := strconv.Atoi(id)
	if err != nil {
		slog.Error("convert string to int error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	var req dto.PaymentRequest

	err = c.ShouldBindJSON(&req)
	if err != nil {
		slog.Error("failed to parse request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	ctx := c.Request.Context()
	payment, err := p.svc.CreatePayment(ctx, parsedId, req)
	if err != nil {
		slog.Error("failed to create payment", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Payment": payment})

}

// GetPaymentById возвращает платеж по id
// @Summary Получить платеж по id
// @Tags payments
// @Produce json
// @Param id path int true "ID платежа"
// @Success 200 {object} map[string]domain.Payment
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /payments/{id} [get]
func (p *PaymentHandlerImpl) GetPaymentById(c *gin.Context) {
	id := c.Param("id")
	parsedId, err := strconv.Atoi(id)
	if err != nil {
		slog.Error("convert string to int error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	ctx := c.Request.Context()
	payment, err := p.svc.GetPaymentById(ctx, parsedId)
	if err != nil {
		slog.Error("failed to get payment by id", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Payment": payment})
}

// GetAllPayments возвращает все платежи
// @Summary Получить список платежей
// @Tags payments
// @Produce json
// @Success 200 {object} map[string][]domain.Payment
// @Failure 500 {object} map[string]string
// @Router /payments [get]
func (p *PaymentHandlerImpl) GetAllPayments(c *gin.Context) {
	ctx := c.Request.Context()
	payments, err := p.svc.GetAllPayments(ctx)
	if err != nil {
		slog.Error("failed to get all payments", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Payments": payments})
}

// UpdatePayment обновляет платеж
// @Summary Обновить платеж
// @Tags payments
// @Accept json
// @Produce json
// @Param id path int true "ID платежа"
// @Param request body dto.UpdatePaymentRequest true "Обновляемые поля"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /payments/{id} [put]
func (p *PaymentHandlerImpl) UpdatePayment(c *gin.Context) {
	id := c.Param("id")
	parsedId, err := strconv.Atoi(id)
	if err != nil {
		slog.Error("convert string to int error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	ctx := c.Request.Context()
	var req dto.UpdatePaymentRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		slog.Error("failed to parse request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	err = p.svc.UpdatePayment(ctx, parsedId, &req)
	if err != nil {
		slog.Error("failed to update payment", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// DeletePayment удаляет платеж
// @Summary Удалить платеж
// @Tags payments
// @Produce json
// @Param id path int true "ID платежа"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /payments/{id} [delete]
func (p *PaymentHandlerImpl) DeletePayment(c *gin.Context) {
	id := c.Param("id")
	parsedId, err := strconv.Atoi(id)
	if err != nil {
		slog.Error("convert string to int error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	ctx := c.Request.Context()
	err = p.svc.DeletePayment(ctx, parsedId)
	if err != nil {
		slog.Error("failed to delete payment", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
