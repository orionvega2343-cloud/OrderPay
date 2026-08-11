package handler

import (
	"OrderPay/internal/order/domain"
	"OrderPay/internal/order/dto"
	"OrderPay/internal/order/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// domain используется только как тип ответа в swagger-аннотациях
var _ domain.Order

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

// PostOrder создает заказ
// @Summary Создать заказ
// @Description Создает новый заказ со списком позиций, считает сумму и выставляет начальный статус
// @Tags orders
// @Accept json
// @Produce json
// @Param request body dto.OrderRequest true "Данные заказа"
// @Success 200 {object} dto.OrderResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders [post]
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

// GetOrderByID возвращает заказ по id
// @Summary Получить заказ по id
// @Tags orders
// @Produce json
// @Param id path int true "ID заказа"
// @Success 200 {object} domain.Order
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders/{id} [get]
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

// GetOrders возвращает все заказы
// @Summary Получить список заказов
// @Tags orders
// @Produce json
// @Success 200 {array} domain.Order
// @Failure 500 {object} map[string]string
// @Router /orders [get]
func (h *OrderHandlerImpl) GetOrders(c *gin.Context) {
	ctx := c.Request.Context()

	order, err := h.svc.GetAllOrders(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}

// UpdateOrder обновляет статус заказа
// @Summary Обновить статус заказа
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "ID заказа"
// @Param request body dto.UpdateOrderStatusRequest true "Новый статус"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders/{id} [patch]
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

// DeleteOrder удаляет заказ
// @Summary Удалить заказ
// @Tags orders
// @Produce json
// @Param id path int true "ID заказа"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders/{id} [delete]
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
