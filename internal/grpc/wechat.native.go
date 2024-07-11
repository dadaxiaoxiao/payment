package grpc

import (
	"context"
	pmtv1 "github.com/dadaxiaoxiao/api-repository/api/proto/gen/payment/v1"
	"github.com/dadaxiaoxiao/payment/internal/domain"
	"github.com/dadaxiaoxiao/payment/internal/service/wechat"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
)

// WechatServiceServer
//
// 微信支付server
type WechatServiceServer struct {
	pmtv1.UnimplementedWechatPaymentServiceServer
	svc *wechat.NativePaymentService
}

func NewWechatServiceServer(svc *wechat.NativePaymentService) *WechatServiceServer {
	return &WechatServiceServer{
		svc: svc,
	}
}

// Register grpc 注册
func (s *WechatServiceServer) Register(server *grpc.Server) {
	pmtv1.RegisterWechatPaymentServiceServer(server, s)
}

// NativePrePay Native 支付
// 如果有H5 支付，使用 H5PrePay
func (s *WechatServiceServer) NativePrePay(ctx context.Context, request *pmtv1.PrePayRequest) (*pmtv1.NativePrePayResponse, error) {
	ctx, span := otel.Tracer("gitee.com/yeqinyiyi/api-repository/payment/service").Start(ctx, "NativePrePay")
	defer func() {
		span.End()
	}()

	codeUrl, err := s.svc.Prepay(ctx, domain.Payment{
		Amt: domain.Amount{
			Total:    request.Amt.GetTotal(),
			Currency: request.Amt.GetCurrency(),
		},
		BizTradeNo:  request.GetBizTradeNo(),
		Description: request.GetDescription(),
	})
	if err != nil {
		return nil, err
	}
	return &pmtv1.NativePrePayResponse{
		CodeUrl: codeUrl,
	}, nil
}

func (s *WechatServiceServer) GetPayment(ctx context.Context, request *pmtv1.GetPaymentRequest) (*pmtv1.GetPaymentResponse, error) {
	ctx, span := otel.Tracer("gitee.com/yeqinyiyi/api-repository/payment/service").Start(ctx, "GetPayment")
	defer func() {
		span.End()
	}()
	p, err := s.svc.GetPayment(ctx, request.GetBizTradeNo())
	if err != nil {
		return nil, err
	}
	return &pmtv1.GetPaymentResponse{
		Status: pmtv1.PaymentStatus((p.Status)),
	}, nil
}
