package handler

import (
	"OrderPay/internal/payment/dto"
	"OrderPay/internal/payment/service"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
