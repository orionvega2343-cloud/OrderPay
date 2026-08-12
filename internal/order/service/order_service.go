package service

import (
	"OrderPay/internal/order/domain"
	"OrderPay/internal/order/dto"
	"OrderPay/internal/order/repository"
	"OrderPay/pkg/transaction"
	"context"
	"log/slog"
)

type OrderService interface {
	CreateOrder(ctx context.Context, req dto.OrderRequest) (*dto.OrderResponse, error)
	GetOrderByID(ctx context.Context, id int) (domain.Order, error)
	GetAllOrders(ctx context.Context) ([]domain.Order, error)
	UpdateOrder(ctx context.Context, id int, req dto.UpdateOrderStatusRequest) error
	DeleteOrder(ctx context.Context, id int) error
}

type OrderServiceImpl struct {
	repo       repository.OrderRepo
	transactor *transaction.Transactor
}

func NewOrderService(repo repository.OrderRepo, transactor *transaction.Transactor) *OrderServiceImpl {
	return &OrderServiceImpl{repo: repo, transactor: transactor}
}

// Создание заказа, устанавливаем начальный статус и UserId,
//в цикле на каждой итерации обходим dto.OrderRequest,
//считаем сумму и через slice записываем поля в domain.OrderItem,
//т.к он вложенный домен в наш заказ,
//оборачиваем repo в transaction для кросс-атомарности нашего метода,
//собираем и возвращаем dto.OrderResponse

func (o *OrderServiceImpl) CreateOrder(ctx context.Context, req dto.OrderRequest) (*dto.OrderResponse, error) {
	m := domain.Order{}
	m.Status = "created"
	m.UserId = req.UserId
	var items []domain.OrderItem
	//Проходим по dto циклом и подставляем поля из JSON ответа
	//вычисляем сумму
	for _, item := range req.Items {
		m.TotalAmount += item.Quantity * item.PricePerUnit
		items = append(items, domain.OrderItem{ProductName: item.ProductName, PricePerUnit: item.PricePerUnit, Quantity: item.Quantity})
	}

	err := o.transactor.WithinTransaction(ctx, func(ctxTx context.Context) error {
		_, err := o.repo.CreateOrder(ctxTx, &m, items)
		slog.Error("failed creating order", "error", err)
		return err
	})
	if err != nil {
		slog.Error("failed check for other errors", "error", err)
		return nil, err
	}
	// TODO: затем добавить конвертацию items → []dto.OrderItemResponse в Service.Create
	var responseItems []dto.OrderItemResponse
	for _, item := range items {
		responseItems = append(responseItems, dto.OrderItemResponse{Id: item.Id, OrderId: item.OrderId, ProductName: item.ProductName, Quantity: item.Quantity, PricePerUnit: item.PricePerUnit})
	}
	return &dto.OrderResponse{Id: m.Id, UserId: m.UserId, Status: m.Status, TotalAmount: m.TotalAmount, Items: responseItems}, nil
}

func (o *OrderServiceImpl) GetOrderByID(ctx context.Context, id int) (domain.Order, error) {
	order, err := o.repo.GetOrderByID(ctx, id)
	if err != nil {
		slog.Error("failed getting order", "error", err)
		return domain.Order{}, err
	}
	return order, nil
}

func (o *OrderServiceImpl) GetAllOrders(ctx context.Context) ([]domain.Order, error) {
	orders, err := o.repo.GetAllOrders(ctx)
	if err != nil {
		slog.Error("failed getting all orders", "error", err)
		return nil, err
	}
	return orders, nil
}

// UpdateOrder - читаем текущий заказ из БД по id,
// вызываем TransitionStatus на прочитанном заказе, передавая
// желаемый статус (req.Status) как новый - метод сам проверяет
// допустимость перехода и меняет статус при успехе,
// сохраняем уже изменённый заказ обратно в БД

func (o *OrderServiceImpl) UpdateOrder(ctx context.Context, id int, req dto.UpdateOrderStatusRequest) error {
	order, err := o.repo.GetOrderByID(ctx, id)
	if err != nil {
		slog.Error("failed getting order", "error", err)
		return err
	}
	err = order.TransitionStatus(req.Status)
	if err != nil {
		slog.Error("failed transitioning status", "error", err)
		return err
	}
	err = o.repo.UpdateOrder(ctx, &order)
	if err != nil {
		slog.Error("failed updating order", "error", err)
		return err
	}
	return nil
}

func (o *OrderServiceImpl) DeleteOrder(ctx context.Context, id int) error {
	err := o.repo.DeleteOrder(ctx, id)
	if err != nil {
		slog.Error("failed delete order", "error", err)
		return err
	}
	return nil
}
