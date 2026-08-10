package service

import (
	"OrderPay/internal/payment/domain"
	"OrderPay/internal/payment/dto"
	"OrderPay/internal/payment/repository"
	"context"
	"errors"
	"log/slog"
)

type RefundService interface {
	CreateRefund(ctx context.Context, paymentId int, req *dto.RefundRequest) (*dto.RefundResponse, error)
	GetAllRefund(ctx context.Context) ([]domain.Refund, error)
}

type RefundServiceImpl struct {
	paymentRepo repository.PaymentRepo
	refundRepo  repository.RefundRepo
}

func NewRefundService(paymentRepo repository.PaymentRepo, refundRepo repository.RefundRepo) *RefundServiceImpl {
	return &RefundServiceImpl{paymentRepo: paymentRepo, refundRepo: refundRepo}
}

// CreateRefund - создает возврат средств пользователю, получает платеж из Payment,
// проверяет что статус обязательно succeeded, валидирует что сумма возврата
// не превышает сумму платежа, создаёт Refund со статусом pending
func (s *RefundServiceImpl) CreateRefund(ctx context.Context, paymentId int, req *dto.RefundRequest) (*dto.RefundResponse, error) {
	p, err := s.paymentRepo.GetPaymentById(ctx, paymentId)
	if err != nil {
		slog.Error("failed getting payment by id", "error", err)
		return nil, err
	}

	if p.Status != "succeeded" {
		slog.Info("payment already created", "paymentId", paymentId)
		return nil, errors.New("payment already created")
	}

	if req.Amount > p.Amount {
		slog.Info("refund amount exceeds payment amount", "amount", req.Amount)
		return nil, errors.New("refund amount exceeds payment amount")
	}

	m := domain.Refund{}
	m.PaymentId = paymentId
	m.Status = "pending"
	m.Reason = req.Reason
	m.Amount = req.Amount
	_, err = s.refundRepo.CreateRefund(ctx, &m)
	if err != nil {
		slog.Error("failed creating refund", "error", err)
		return nil, err
	}
	return &dto.RefundResponse{Id: m.Id, PaymentId: m.PaymentId, Amount: m.Amount, Reason: m.Reason, Status: m.Status}, nil
}

func (s *RefundServiceImpl) GetAllRefund(ctx context.Context) ([]domain.Refund, error) {
	refunds, err := s.refundRepo.GetAllRefund(ctx)
	if err != nil {
		slog.Error("failed getting all refunds", "error", err)
		return nil, err
	}
	return refunds, nil
}
