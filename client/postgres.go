package client

import (
	"fmt"
	"log"

	"github.com/shashank-mugiwara/laughingtale/conf"
	"github.com/shashank-mugiwara/laughingtale/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var postgresDb *gorm.DB

func GetPostgresDb() *gorm.DB {
	return postgresDb
}

func InitPostgresPoller() {

	if conf.PostgresSetting.Enabled == "false" {
		return
	}

	dns := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", conf.PostgresSetting.ReaderNode, conf.PostgresSetting.Username,
		conf.PostgresSetting.Password, conf.PostgresSetting.Database, conf.PostgresSetting.Port, conf.PostgresSetting.SslEnabled)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dns,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		log.Fatalln("Unable to connect to Postgres DB. Plese check if config is correct and DB is running")
	}

	logger.GetLaughingTaleLogger().Info("Connection to PostgresDb is successful.")
	postgresDb = db
}
