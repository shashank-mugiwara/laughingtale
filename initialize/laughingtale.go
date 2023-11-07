package initialize

import (
	"github.com/shashank-mugiwara/laughingtale/client"
	"github.com/shashank-mugiwara/laughingtale/conf"
	customLogger "github.com/shashank-mugiwara/laughingtale/logger"
	"github.com/shashank-mugiwara/laughingtale/pkg/connectors"
	"github.com/shashank-mugiwara/laughingtale/pkg/metrics"
)

func InitRoutes() {
	metrics.RegisterRoutes(conf.GetLaughingTaleEngine(), customLogger.GetLaughingTaleLogger())
	connectors.RegisterRoutes(conf.GetLaughingTaleEngine(), customLogger.GetLaughingTaleLogger())
}

func InitConfig() {
	conf.SetUp("conf/config.ini")
}

func InitKafkaConsumer() {
	client.InitKafkaConsumer()
}
