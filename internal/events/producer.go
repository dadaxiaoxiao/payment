package events

import "context"

type Producer interface {
	ProducerPaymentEvent(cxt context.Context, evt PaymentEvent) error
}
