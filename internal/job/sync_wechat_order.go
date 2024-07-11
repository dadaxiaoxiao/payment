package job

import (
	"context"
	"github.com/dadaxiaoxiao/go-pkg/accesslog"
	"github.com/dadaxiaoxiao/payment/internal/service/wechat"
	"time"
)

type SyncWechatOrderJob struct {
	log accesslog.Logger
	svc *wechat.NativePaymentService
}

func (s *SyncWechatOrderJob) Name() string {
	return "sync_wechat_order_job"
}

func (s *SyncWechatOrderJob) Run() error {
	offset := 0
	// 也可以做成参数
	const limit = 100
	// 三十分钟之前的订单我们就认为已经过期了。
	now := time.Now().Add(-time.Minute * 30)
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		pmts, err := s.svc.FindExpiredPayment(ctx, offset, limit, now)
		// 这里的超时控制
		cancel()
		if err != nil {
			// 直接中断
			return err
		}
		for _, pmt := range pmts {
			// 单个重新设置超时
			ctx, cancel = context.WithTimeout(context.Background(), time.Second)
			err = s.svc.SyncWechatInfo(ctx, pmt.BizTradeNo)
			if err != nil {
				s.log.Error("同步微信支付信息失败",
					accesslog.String("trade_no", pmt.BizTradeNo),
					accesslog.Error(err))
			}
			cancel()
		}
		if len(pmts) < limit {
			// 没数据了
			return nil
		}
		offset = offset + len(pmts)
	}
}

func NewSyncWechatOrderJob(log accesslog.Logger, svc *wechat.NativePaymentService) *SyncWechatOrderJob {
	return &SyncWechatOrderJob{
		log: log,
		svc: svc,
	}
}
