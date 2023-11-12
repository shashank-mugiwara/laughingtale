package client

import (
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

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=localhost user=shashank.j password=root dbname=maindb port=5432 sslmode=disable",
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		log.Fatalln("Unable to connect to Postgres DB. Plese check if config is correct and DB is running")
	}

	logger.GetLaughingTaleLogger().Info("Connection to PostgresDb is successful.")
	postgresDb = db
}
