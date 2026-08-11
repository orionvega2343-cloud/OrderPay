package handler

import (
	"OrderPay/internal/payment/service"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AggregationHandler interface {
	GetPaymentTotalAmount(c *gin.Context)
}
type AggregationHandlerImpl struct {
	svc service.AggregationService
}

func NewAggregationHandlerImpl(svc service.AggregationService) *AggregationHandlerImpl {
	return &AggregationHandlerImpl{svc: svc}
}

func (h *AggregationHandlerImpl) GetPaymentTotalAmount(c *gin.Context) {
	ctx := c.Request.Context()
	userId := c.Query("user_id")
	status := c.Query("status")
	from := c.Query("from")
	to := c.Query("to")
	layout := "2006-01-02 15:04"

	parsedFrom, err := time.Parse(layout, from)
	if err != nil {
		slog.Error("failed to parse from", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	parsedTo, err := time.Parse(layout, to)
	if err != nil {
		slog.Error("failed to parse to", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	a, err := h.svc.GetPaymentTotalAmount(ctx, userId, status, parsedFrom, parsedTo)
	if err != nil {
		slog.Error("failed to get payment total amount", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{"aggregation": a})
}
