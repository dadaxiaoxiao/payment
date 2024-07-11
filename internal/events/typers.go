package events

// PaymentEvent
// 支付事件信息
type PaymentEvent struct {
	BizTradeNo string
	Status     uint8
}

// Topic 获取kafka topic
func (PaymentEvent) Topic() string {
	// 这里暂时统一使用同一个topic
	// 也可以不同业务支付，使用不同的 biz 进行 biz + "_payment_events"
	return "payment_events"
}
