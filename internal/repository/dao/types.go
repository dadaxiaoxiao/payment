package dao

import (
	"context"
	"database/sql"
	"github.com/dadaxiaoxiao/payment/internal/domain"
	"time"
)

type PaymentDao interface {
	Insert(ctx context.Context, pmt Payment) error
	UpdateTxnIDAndStatus(ctx context.Context, bizTradeNo string, txnID string, status domain.PaymentStatus) error
	FindExpiredPayment(ctx context.Context, offset int, limit int, t time.Time) ([]Payment, error)
	GetPayment(ctx context.Context, bizTradeNo string) (Payment, error)
}

type Payment struct {
	Id int64 `gorm:"primaryKey,autoIncrement" bson:"id,omitempty"`

	// 订单金额
	Amt int64 `gorm:"column:amt"`

	// 货币类型
	Currency string `gorm:"column:currency;type:varchar(128);"`

	// 简短的描述
	Description string `gorm:"column:description"`

	// 业务方传过来的
	BizTradeNO string `gorm:"column:biz_trade_no;type:varchar(256);unique"`

	// 第三方支付平台的事务 ID，唯一的
	TxnID sql.NullString `gorm:"column:txn_id;type:varchar(128);unique"`

	Status uint8 `gorm:"column:status"`
	Utime  int64
	Ctime  int64
}
