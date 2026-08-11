package handler

import (
	"OrderPay/internal/payment/dto"
	"OrderPay/internal/payment/service"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RefundHandler interface {
	CreateRefund(c *gin.Context)
	GetAllRefund(c *gin.Context)
}

type RefundHandlerImpl struct {
	svc service.RefundService
}

func NewRefundHandler(svc service.RefundService) *RefundHandlerImpl {
	return &RefundHandlerImpl{svc: svc}
}

func (h *RefundHandlerImpl) CreateRefund(c *gin.Context) {
	id := c.Param("id")
	parsedId, err := strconv.Atoi(id)
	if err != nil {
		slog.Error("Parse id failed", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	var req dto.RefundRequest

	err = c.ShouldBindJSON(&req)
	if err != nil {
		slog.Error("Parse req failed", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	ctx := c.Request.Context()
	refund, err := h.svc.CreateRefund(ctx, parsedId, &req)
	if err != nil {
		slog.Error("CreateRefund failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{"refund": refund})
}

func (h *RefundHandlerImpl) GetAllRefund(c *gin.Context) {
	ctx := c.Request.Context()
	refunds, err := h.svc.GetAllRefund(ctx)
	if err != nil {
		slog.Error("GetAllRefund failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{"refunds": refunds})
}