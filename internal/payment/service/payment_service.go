package service

import (
	orderRepo "OrderPay/internal/order/repository"
	"OrderPay/internal/payment/domain"
	"OrderPay/internal/payment/dto"
	"OrderPay/internal/payment/repository"
	"OrderPay/pkg/transaction"
	"context"
	"errors"
	"log/slog"
)

type PaymentService interface {
	CreatePayment(ctx context.Context, orderId int, req dto.PaymentRequest) (dto.PaymentResponse, error)
	GetPaymentById(ctx context.Context, id int) (*domain.Payment, error)
	GetAllPayments(ctx context.Context) ([]domain.Payment, error)
	UpdatePayment(ctx context.Context, req *dto.UpdatePaymentRequest) error
	DeletePayment(ctx context.Context, id int) error
}

type PaymentServiceImpl struct {
	repo       repository.PaymentRepo
	orderRepo  orderRepo.OrderRepo
	transactor *transaction.Transactor
}

func NewPaymentService(repo repository.PaymentRepo, transactor *transaction.Transactor, orderRepo orderRepo.OrderRepo) *PaymentServiceImpl {
	return &PaymentServiceImpl{repo: repo, transactor: transactor, orderRepo: orderRepo}
}

// CreatePayment - получает заказ по id, проверяет что его статус created,
// собирает Payment с суммой из заказа (сумма не приходит от клиента),
// переводит статус заказа в paid через TransitionStatus (в памяти, до транзакции),
// оборачивает создание Payment и сохранение обновлённого заказа в одну
// транзакцию для кросс-доменной атомарности

func (s *PaymentServiceImpl) CreatePayment(ctx context.Context, orderId int, req dto.PaymentRequest) (dto.PaymentResponse, error) {
	order, err := s.orderRepo.GetOrderByID(ctx, orderId)
	if err != nil {
		slog.Error("failed getting order repository", "error", err)
		return dto.PaymentResponse{}, err
	}

	if order.Status != "created" {
		slog.Error("failed getting order repository", "error", order.Status)
		return dto.PaymentResponse{}, errors.New("order is not in created status")
	}

	m := domain.Payment{}
	m.OrderId = orderId
	m.Amount = order.TotalAmount
	m.Method = req.Method
	m.Status = "succeeded"

	err = order.TransitionStatus("paid")
	if err != nil {
		slog.Error("failed to transition order status", "error", err)
		return dto.PaymentResponse{}, errors.New("failed updating order status")
	}

	err = s.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		_, err := s.repo.CreatePayment(ctx, &m)
		if err != nil {
			slog.Error("failed creating payment repository", "error", err)
			return err
		}
		err = s.orderRepo.UpdateOrder(ctx, &order)
		if err != nil {
			slog.Error("failed updating order repository", "error", err)
			return err
		}
		return nil
	})
	if err != nil {
		slog.Error("failed creating payment repository", "error", err)
		return dto.PaymentResponse{}, err
	}
	return dto.PaymentResponse{Id: m.Id, OrderId: m.OrderId, Amount: m.Amount, Status: m.Status, Method: m.Method}, nil
}

func (s *PaymentServiceImpl) GetPaymentById(ctx context.Context, id int) (*domain.Payment, error) {
	p, err := s.repo.GetPaymentById(ctx, id)
	if err != nil {
		slog.Error("failed getting payment repository", "error", err)
		return nil, err
	}
	return p, nil
}

func (s *PaymentServiceImpl) GetAllPayments(ctx context.Context) ([]domain.Payment, error) {
	ps, err := s.repo.GetAllPayments(ctx)
	if err != nil {
		slog.Error("failed getting payment repository", "error", err)
		return nil, err
	}
	return ps, nil
}

func (s *PaymentServiceImpl) UpdatePayment(ctx context.Context, id int, req *dto.UpdatePaymentRequest) error {
	m := domain.Payment{}
	m.Id = id
	if req.Status != "pending" && req.Status != "succeeded" && req.Status != "failed" {
		return errors.New("invalid status request")
	}
	m.Status = req.Status
	m.Method = req.Method
	err := s.repo.UpdatePayment(ctx, &m)
	if err != nil {
		slog.Error("Failed updating payment repository", "error", err)
		return err
	}
	return nil
}

func (s *PaymentServiceImpl) DeletePayment(ctx context.Context, id int) error {
	err := s.repo.DeletePayment(ctx, id)
	if err != nil {
		slog.Error("failed deleting payment repository", "error", err)
		return err
	}
	return nil
}
