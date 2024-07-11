package ioc

import (
	"context"
	"github.com/dadaxiaoxiao/go-pkg/accesslog"
	"github.com/dadaxiaoxiao/payment/internal/events"
	"github.com/dadaxiaoxiao/payment/internal/repository"
	"github.com/dadaxiaoxiao/payment/internal/service/wechat"
	"github.com/spf13/viper"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

func InitWechatClient(cfg WechatConfig) *core.Client {
	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(cfg.KeyPath)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(cfg.MchID, cfg.MchSerialNum, mchPrivateKey, cfg.MchAPIv3Key),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		panic(err)
	}
	return client
}

func InitWechatNativePaymentSvc(
	cli *core.Client,
	repo repository.PaymentRepository,
	l accesslog.Logger,
	cfg WechatConfig,
	producer events.Producer) *wechat.NativePaymentService {
	return wechat.NewNativePaymentService(
		&native.NativeApiService{Client: cli},
		repo,
		l, producer,
		cfg.AppID, cfg.MchID)
}

func InitWechatNotifyHandler(cfg WechatConfig) *notify.Handler {
	certificateVisitor := downloader.MgrInstance().GetCertificateVisitor(cfg.MchID)
	// 3. 使用apiv3 key、证书访问器初始化 `notify.Handler`
	handler, err := notify.NewRSANotifyHandler(cfg.MchAPIv3Key,
		verifiers.NewSHA256WithRSAVerifier(certificateVisitor))
	if err != nil {
		panic(err)
	}
	return handler
}

// InitWechatConfig
//
// 初始化 wechatConfig 这里可以考虑从环境变变量，或者配置中心获取
func InitWechatConfig() WechatConfig {

	var config WechatConfig
	err := viper.UnmarshalKey("wechatConfig", &config)
	if err != nil {
		panic(err)
	}
	return config
}

type WechatConfig struct {
	// 小程序id
	AppID string `yaml:"appID"`
	// 商户id
	MchID string `yaml:"mchID"`
	// 商户APIv3密钥
	MchAPIv3Key string `yaml:"mchAPIv3Key"`
	//商户API 证书的序列号
	MchSerialNum string `yaml:"mchSerialNum"`

	// 证书
	CertPath string `yaml:"certPath"`
	//商户API私钥
	KeyPath string `yaml:"keyPath"`
}
