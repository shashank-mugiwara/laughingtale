package initialize

import (
	"github.com/shashank-mugiwara/laughingtale/conf"
	customLogger "github.com/shashank-mugiwara/laughingtale/logger"
	"github.com/shashank-mugiwara/laughingtale/metrics"
)

func InitRoutes() {
	metrics.RegisterRoutes(conf.GetLaughingTaleEngine(), customLogger.GetLaughingTaleLogger())
}

func InitConfig() {
	conf.SetUp("conf/config.ini")
}
