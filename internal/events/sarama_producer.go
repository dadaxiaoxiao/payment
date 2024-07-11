package events

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
)

// SaramaProducer
//
// 通过 Sarama 实现了 Producer的接口
type SaramaProducer struct {
	producer sarama.SyncProducer
}

// NewSaramaProducer 新建SaramaProducer
//
// 这里只需要注入 sarama.Client
func NewSaramaProducer(client sarama.Client) (*SaramaProducer, error) {
	p, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return nil, err
	}
	return &SaramaProducer{
		producer: p,
	}, nil
}

func (s *SaramaProducer) ProducerPaymentEvent(cxt context.Context, evt PaymentEvent) error {
	data, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	_, _, err = s.producer.SendMessage(&sarama.ProducerMessage{
		Key:   sarama.StringEncoder(evt.BizTradeNo),
		Topic: evt.Topic(),
		Value: sarama.ByteEncoder(data),
	})
	return err
}
