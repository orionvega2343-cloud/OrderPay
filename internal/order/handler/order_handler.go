package handler

import (
	"OrderPay/internal/order/dto"
	"OrderPay/internal/order/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderHandler interface {
	PostOrder(c *gin.Context)
	GetOrderByID(c *gin.Context)
	GetOrders(c *gin.Context)
	UpdateOrder(c *gin.Context)
	DeleteOrder(c *gin.Context)
}

type OrderHandlerImpl struct {
	svc service.OrderService
}

func NewOrderHandlerImpl(svc service.OrderService) *OrderHandlerImpl {
	return &OrderHandlerImpl{svc: svc}
}

// TODO: sentinel errors (например, ErrInvalidTransition в domain-пакете) +
// errors.Is() в хендлере, чтобы отличать ошибки валидации перехода статуса (400)
// от внутренних ошибок сервера (500) - сейчас все ошибки уходят как 500

func (h *OrderHandlerImpl) PostOrder(c *gin.Context) {
	var req dto.OrderRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	order, err := h.svc.CreateOrder(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandlerImpl) GetOrderByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	parsedId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.svc.GetOrderByID(ctx, parsedId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *OrderHandlerImpl) GetOrders(c *gin.Context) {
	ctx := c.Request.Context()

	order, err := h.svc.GetAllOrders(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *OrderHandlerImpl) UpdateOrder(c *gin.Context) {
	id := c.Param("id")
	parsedId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req dto.UpdateOrderStatusRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	err = h.svc.UpdateOrder(ctx, parsedId, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (h *OrderHandlerImpl) DeleteOrder(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	parsedId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.svc.DeleteOrder(ctx, parsedId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
