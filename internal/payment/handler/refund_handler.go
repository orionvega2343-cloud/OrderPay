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
var _ domain.Refund

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

// CreateRefund создает возврат по платежу
// @Summary Создать возврат
// @Description id в пути - это id платежа, по которому оформляется возврат
// @Tags refunds
// @Accept json
// @Produce json
// @Param id path int true "ID платежа"
// @Param request body dto.RefundRequest true "Данные возврата"
// @Success 200 {object} map[string]dto.RefundResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /payments/{id}/refunds [post]
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

// GetAllRefund возвращает все возвраты
// @Summary Получить список возвратов
// @Tags refunds
// @Produce json
// @Success 200 {object} map[string][]domain.Refund
// @Failure 500 {object} map[string]string
// @Router /refunds [get]
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