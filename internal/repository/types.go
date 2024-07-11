package repository

import (
	"context"
	"github.com/dadaxiaoxiao/payment/internal/domain"
	"time"
)

type PaymentRepository interface {
	// AddPayment 添加支付
	AddPayment(ctx context.Context, payment domain.Payment) error
	// UpdatePayment 修改支付 （一般只用来修改支付状态）
	UpdatePayment(ctx context.Context, payment domain.Payment) error
	// FindExpiredPayment 查找过期支付
	FindExpiredPayment(ctx context.Context, offset int, limit int, t time.Time) ([]domain.Payment, error)
	// GetPayment 获取支付信息
	GetPayment(ctx context.Context, bizTradeNo string) (domain.Payment, error)
}
