package service

import (
	"OrderPay/internal/payment/repository"
	"context"
	"time"
)

type AggregationService interface {
	GetPaymentTotalAmount(ctx context.Context, userId, status string, from, to time.Time) (int, error)
}

type AggregationServiceImpl struct {
	repo repository.PaymentRepo
}

func NewAggregationService(repo repository.PaymentRepo) *AggregationServiceImpl {
	return &AggregationServiceImpl{repo: repo}
}

func (s *AggregationServiceImpl) GetPaymentTotalAmount(ctx context.Context, userId, status string, from, to time.Time) (int, error) {
	total, err := s.repo.GetPaymentTotalAmount(ctx, userId, status, from, to)
	if err != nil {
		return 0, err
	}
	return total, nil
}
