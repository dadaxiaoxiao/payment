package wechat

import (
	"context"
	"errors"
	"fmt"
	"github.com/dadaxiaoxiao/go-pkg/accesslog"
	"github.com/dadaxiaoxiao/payment/internal/domain"
	"github.com/dadaxiaoxiao/payment/internal/events"
	"github.com/dadaxiaoxiao/payment/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"time"
)

var errUnknownTransactionState = errors.New("未知的微信事务状态")

type NativePaymentService struct {
	svc *native.NativeApiService
	// 小程序id
	appID string
	// 商户号id
	mchID     string
	notifyURL string
	repo      repository.PaymentRepository
	l         accesslog.Logger

	// wechat 交易状态 与 domain.PaymentStatus 映射关系
	// SUCCESS：支付成功
	// REFUND：转入退款
	// NOTPAY：未支付
	// CLOSED：已关闭
	// REVOKED：已撤销（仅付款码支付会返回）
	// USERPAYING：用户支付中（仅付款码支付会返回）
	// PAYERROR：支付失败（仅付款码支付会返回）
	nativeCBTypeToStatus map[string]domain.PaymentStatus
	producer             events.Producer
}

func NewNativePaymentService(svc *native.NativeApiService,
	repo repository.PaymentRepository,
	l accesslog.Logger,
	producer events.Producer,
	appID string, mchID string) *NativePaymentService {
	return &NativePaymentService{
		svc:      svc,
		repo:     repo,
		l:        l,
		producer: producer,
		appID:    appID,
		mchID:    mchID,
		// 回调url 一般来说是固定的
		notifyURL: "http://qinyeyiyi.com/pay/callback",
		nativeCBTypeToStatus: map[string]domain.PaymentStatus{
			"SUCCESS":    domain.PaymentStatusSuccess,
			"PAYERROR":   domain.PaymentStatusFailed,
			"NOTPAY":     domain.PaymentStatusInit,
			"CLOSED":     domain.PaymentStatusFailed,
			"REVOKED":    domain.PaymentStatusFailed,
			"REFUND":     domain.PaymentStatusRefund,
			"USERPAYING": domain.PaymentStatusInit,
		},
	}
}

func (n *NativePaymentService) Prepay(ctx context.Context, pmt domain.Payment) (string, error) {
	// 唯一索引冲突
	// 业务方唤起了支付，但是没付，下一次再过来，应该换 BizTradeNO
	err := n.repo.AddPayment(ctx, pmt)
	if err != nil {
		return "", err
	}

	resp, result, err := n.svc.Prepay(ctx,
		native.PrepayRequest{
			Appid:       core.String(n.appID),
			Mchid:       core.String(n.mchID),
			Description: core.String(pmt.Description),
			// 这里商户订单号是业务方传递过来，不主动生成，这里一定要一个唯一标识
			OutTradeNo: core.String(pmt.BizTradeNo),
			NotifyUrl:  core.String(n.notifyURL),
			// 这里设置半小时有效期
			TimeExpire: core.Time(time.Now().Add(time.Minute * 30)),
			Amount: &native.Amount{
				Total:    core.Int64(pmt.Amt.Total),
				Currency: core.String(pmt.Amt.Currency),
			},
		})

	n.l.Debug("微信prepay响应",
		accesslog.Field{Key: "result", Value: result},
		accesslog.Field{Key: "resp", Value: resp})
	if err != nil {
		return "", err
	}
	return *resp.CodeUrl, nil
}

// HandleCallback
//
// 回调句柄
func (n *NativePaymentService) HandleCallback(ctx *gin.Context, txn *payments.Transaction) error {
	return n.updateByTxn(ctx, txn)
}

// FindExpiredPayment
//
// 查找过期的支付
func (n *NativePaymentService) FindExpiredPayment(ctx context.Context, offset, limit int, t time.Time) ([]domain.Payment, error) {
	return n.repo.FindExpiredPayment(ctx, offset, limit, t)
}

// SyncWechatInfo 同步微信信息
// 这里用于定时对账任何，请求微信支付信息
// 商户订单号查询订单
func (n *NativePaymentService) SyncWechatInfo(ctx context.Context, bizTradNo string) error {
	resp, result, err := n.svc.QueryOrderByOutTradeNo(ctx, native.QueryOrderByOutTradeNoRequest{
		OutTradeNo: core.String(bizTradNo),
		Mchid:      core.String(n.mchID),
	})
	if err != nil {
		n.l.Error("微信商户订单号查询订单响应",
			accesslog.Field{Key: "result", Value: result},
			accesslog.Field{Key: "resp", Value: resp})
		return err
	}
	return n.updateByTxn(ctx, resp)
}

func (n *NativePaymentService) GetPayment(ctx context.Context, bizTradeId string) (domain.Payment, error) {
	res, err := n.repo.GetPayment(ctx, bizTradeId)
	/*
			这里如果没有知道支付结果，是否要设计慢路径
		    比如 去微信查询，然后更新本地数据
	*/
	return res, err
}

func (n *NativePaymentService) updateByTxn(ctx context.Context, txn *payments.Transaction) error {
	// 微信支付状态转为 自定义状态
	status, ok := n.nativeCBTypeToStatus[*txn.TradeState]
	if !ok {
		return fmt.Errorf("%w, %s", errUnknownTransactionState, *txn.TradeState)
	}
	err := n.repo.UpdatePayment(ctx, domain.Payment{
		BizTradeNo: *txn.OutTradeNo,
		TxnID:      *txn.TransactionId,
		Status:     status,
	})
	if err != nil {
		return err
	}
	// 对账消息写入消息队列
	// 存在问题：
	// 1.微信成功回调，但是本地同步失败 （部分失败问题）
	// 2.如果系统没有返回 200 状态码，会重试，这里会重复发送

	err1 := n.producer.ProducerPaymentEvent(ctx, events.PaymentEvent{
		BizTradeNo: *txn.OutTradeNo,
		Status:     status.AsUint8(),
	})
	if err1 != nil {
		// 加监控加告警
		// 立刻手动修复，或者自动补发
	}
	return nil
}
