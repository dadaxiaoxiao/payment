package web

import (
	"github.com/dadaxiaoxiao/go-pkg/accesslog"
	"github.com/dadaxiaoxiao/go-pkg/ginx"
	"github.com/dadaxiaoxiao/payment/internal/service/wechat"
	"github.com/gin-gonic/gin"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
)

type WechatHandler struct {
	handler *notify.Handler
	log     accesslog.Logger
	// 这没有使用统一支付接口的抽象
	nativeSvc *wechat.NativePaymentService
}

func NewWechatHandler(handler *notify.Handler, nativeSvc *wechat.NativePaymentService, log accesslog.Logger) *WechatHandler {
	return &WechatHandler{
		handler:   handler,
		nativeSvc: nativeSvc,
		log:       log,
	}
}

func (w *WechatHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/wechat")
	ug.Any("/pay/callback", ginx.Wrap(w.HandleNative))
}

func (w *WechatHandler) HandleNative(ctx *gin.Context) (ginx.Result, error) {
	transaction := &payments.Transaction{}
	_, err := w.handler.ParseNotifyRequest(ctx, ctx.Request, transaction)
	if err != nil {
		return ginx.Result{}, err
	}
	err = w.nativeSvc.HandleCallback(ctx, transaction)
	return ginx.Result{}, err
}
