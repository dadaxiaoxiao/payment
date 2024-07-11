//go:build wireinject

package main

import (
	"github.com/dadaxiaoxiao/go-pkg/customserver"
	grpc2 "github.com/dadaxiaoxiao/payment/internal/grpc"
	"github.com/dadaxiaoxiao/payment/internal/repository"
	dao2 "github.com/dadaxiaoxiao/payment/internal/repository/dao"
	"github.com/dadaxiaoxiao/payment/internal/web"
	"github.com/dadaxiaoxiao/payment/ioc"
	"github.com/google/wire"
)

var thirdPartyProvider = wire.NewSet(
	ioc.InitRedis,
	ioc.InitOTEL,
	ioc.InitLogger,
	ioc.InitDB,
	ioc.InitKafka,
	ioc.InitProducer,
	ioc.InitEtcdClient,
)

var jobProvider = wire.NewSet(
	ioc.InitCronJobBuilder,
	ioc.InitJobs,
)

func InitApp() *customserver.App {
	wire.Build(
		thirdPartyProvider,
		jobProvider,
		ioc.InitWechatConfig,
		ioc.InitWechatClient,
		ioc.InitWechatNotifyHandler,
		dao2.NewPaymentGORMDAO,
		repository.NewPaymentRepository,
		ioc.InitWechatNativePaymentSvc,
		grpc2.NewWechatServiceServer,
		web.NewWechatHandler,
		ioc.InitGinServer,
		ioc.InitGRPCServer,
		wire.Struct(new(customserver.App), "GRPCServer", "GinServer", "Crons"))
	return new(customserver.App)
}
