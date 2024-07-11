package ioc

import (
	"github.com/dadaxiaoxiao/go-pkg/ginx"
	"github.com/dadaxiaoxiao/payment/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func InitGinServer(hdl *web.WechatHandler) *ginx.Server {
	engine := gin.Default()
	addr := viper.GetString("http.addr")
	hdl.RegisterRoutes(engine)
	return &ginx.Server{
		Engine: engine,
		Addr:   addr,
	}
}
