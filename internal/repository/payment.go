package repository

import (
	"context"
	"github.com/dadaxiaoxiao/payment/internal/domain"
	"github.com/dadaxiaoxiao/payment/internal/repository/dao"
	"github.com/ecodeclub/ekit/slice"
	"time"
)

type paymentRepository struct {
	dao dao.PaymentDao
}

func NewPaymentRepository(dao dao.PaymentDao) PaymentRepository {
	return &paymentRepository{
		dao: dao,
	}
}

func (p *paymentRepository) AddPayment(ctx context.Context, payment domain.Payment) error {
	return p.dao.Insert(ctx, p.toEntity(payment))
}

func (p *paymentRepository) UpdatePayment(ctx context.Context, payment domain.Payment) error {
	return p.dao.UpdateTxnIDAndStatus(ctx, payment.BizTradeNo, payment.TxnID, payment.Status)
}

func (p *paymentRepository) FindExpiredPayment(ctx context.Context, offset int, limit int, t time.Time) ([]domain.Payment, error) {
	pmts, err := p.dao.FindExpiredPayment(ctx, offset, limit, t)
	if err != nil {
		return nil, err
	}

	res := slice.Map(pmts, func(idx int, src dao.Payment) domain.Payment {
		return p.toDomain(src)
	})
	return res, nil
}

func (p *paymentRepository) GetPayment(ctx context.Context, bizTradeNo string) (domain.Payment, error) {
	data, err := p.dao.GetPayment(ctx, bizTradeNo)
	return p.toDomain(data), err
}

func (p *paymentRepository) toEntity(pmt domain.Payment) dao.Payment {
	return dao.Payment{
		Amt:         pmt.Amt.Total,
		Currency:    pmt.Amt.Currency,
		BizTradeNO:  pmt.BizTradeNo,
		Description: pmt.Description,
		Status:      domain.PaymentStatusInit,
	}
}

func (p *paymentRepository) toDomain(pmt dao.Payment) domain.Payment {
	return domain.Payment{
		Amt: domain.Amount{
			Currency: pmt.Currency,
			Total:    pmt.Amt,
		},
		BizTradeNo:  pmt.BizTradeNO,
		Description: pmt.Description,
		Status:      domain.PaymentStatus(pmt.Status),
		TxnID:       pmt.TxnID.String,
	}
}
