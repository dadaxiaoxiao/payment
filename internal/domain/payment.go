package domain

// Amount 价格
type Amount struct {
	//【货币类型】
	Currency string
	// 【总金额】 订单总金额，单位为分
	Total int64
}

// Payment 支付
type Payment struct {
	Amt Amount
	// 业务订单号
	BizTradeNo string
	// 订单描述
	Description string

	Status PaymentStatus
	// 第三方那边返回的 事务id
	TxnID string
}

// WePayment 微信支付
type WePayment struct {
	Payment
}

// PaymentStatus 支付状态
type PaymentStatus uint8

func (s PaymentStatus) AsUint8() uint8 {
	return uint8(s)
}

const (
	PaymentStatusUnknown = iota
	PaymentStatusInit
	PaymentStatusSuccess
	PaymentStatusFailed
	PaymentStatusRefund
)
