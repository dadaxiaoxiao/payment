package ioc

import (
	"github.com/dadaxiaoxiao/go-pkg/accesslog"
	"github.com/dadaxiaoxiao/go-pkg/jobx"
	"github.com/dadaxiaoxiao/payment/internal/job"
	"github.com/dadaxiaoxiao/payment/internal/service/wechat"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
)

func InitCronJobBuilder(log accesslog.Logger) *jobx.CronJobBuilder {
	cron_job_builder := jobx.NewCronJobBuilder("gitee.com/yeqinyiyi/basic-go/webook/payment", log,
		prometheus.SummaryOpts{
			Namespace: "qinye",
			Subsystem: "demo",
			Help:      "统计定时任务的执行情况",
			Name:      "cron_payment_job",
		})
	return cron_job_builder
}

// InitJobs 初始job
//
// 这里不可避免要依赖注入和业务相关的svc
func InitJobs(log accesslog.Logger,
	jobBuilder *jobx.CronJobBuilder,
	svc *wechat.NativePaymentService,
) []*cron.Cron {
	return []*cron.Cron{
		initSyncWechatOrderJob(log, jobBuilder, svc),
	}
}

func initSyncWechatOrderJob(log accesslog.Logger,
	jobBuilder *jobx.CronJobBuilder,
	svc *wechat.NativePaymentService) *cron.Cron {
	res := cron.New(cron.WithSeconds())
	cron_job := jobBuilder.Build(job.NewSyncWechatOrderJob(log, svc))
	// 这里每三分钟一次
	_, err := res.AddJob("0 */3 * * * ?", cron_job)
	if err != nil {
		panic(err)
	}
	return res
}
