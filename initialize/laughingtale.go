package initialize

import (
	"github.com/shashank-mugiwara/laughingtale/client"
	"github.com/shashank-mugiwara/laughingtale/conf"
	"github.com/shashank-mugiwara/laughingtale/db"
	"github.com/shashank-mugiwara/laughingtale/logger"
	customLogger "github.com/shashank-mugiwara/laughingtale/logger"
	"github.com/shashank-mugiwara/laughingtale/pkg/metrics"
	sourceconfig "github.com/shashank-mugiwara/laughingtale/pkg/source_config"
)

func InitRoutes() {
	metrics.RegisterRoutes(conf.GetLaughingTaleEngine(), customLogger.GetLaughingTaleLogger())
	sourceconfig.RegisterRoutes(conf.GetLaughingTaleEngine(), customLogger.GetLaughingTaleLogger())
}

func InitConfig() {
	conf.SetUp("conf/config.ini")
	logger.InitLaughingTaleLogger()
}

func InitClients() {
	client.InitKafkaConsumer()
	client.InitMongodb()
	db.InitGormPool()
	client.InitRedisClient()
}
