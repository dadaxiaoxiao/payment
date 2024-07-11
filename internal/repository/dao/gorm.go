package dao

import (
	"context"
	"github.com/dadaxiaoxiao/payment/internal/domain"
	"gorm.io/gorm"
	"time"
)

type PaymentGORMDAO struct {
	db *gorm.DB
}

func NewPaymentGORMDAO(db *gorm.DB) PaymentDao {
	return &PaymentGORMDAO{db: db}
}

func (p *PaymentGORMDAO) Insert(ctx context.Context, pmt Payment) error {
	now := time.Now().UnixMilli()
	pmt.Utime = now
	pmt.Ctime = now
	return p.db.WithContext(ctx).Create(&pmt).Error
}

func (p *PaymentGORMDAO) UpdateTxnIDAndStatus(ctx context.Context,
	bizTradeNo string,
	txnID string, status domain.PaymentStatus) error {
	return p.db.WithContext(ctx).
		Where("biz_trade_no =?", bizTradeNo).Model(&Payment{}).
		Updates(map[string]any{
			"txn_id": txnID,
			"status": status.AsUint8(),
			"utime":  time.Now().UnixMilli(),
		}).Error
}

func (p *PaymentGORMDAO) FindExpiredPayment(ctx context.Context, offset int, limit int, t time.Time) ([]Payment, error) {
	var res []Payment
	err := p.db.WithContext(ctx).Where("status = ? AND utime < ?",
		uint8(domain.PaymentStatusInit), t.UnixMilli()).
		Offset(offset).Limit(limit).Find(&res).Error
	return res, err
}

func (p *PaymentGORMDAO) GetPayment(ctx context.Context, bizTradeNo string) (Payment, error) {
	var res Payment
	err := p.db.WithContext(ctx).Where("biz_trade_no = ?", bizTradeNo).First(&res).Error
	return res, err
}
