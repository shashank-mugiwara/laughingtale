package poller

import (
	"context"
	"time"

	"github.com/shashank-mugiwara/laughingtale/client"
	"github.com/shashank-mugiwara/laughingtale/db"
	"github.com/shashank-mugiwara/laughingtale/logger"
	"github.com/shashank-mugiwara/laughingtale/pkg/factory"
	"github.com/shashank-mugiwara/laughingtale/pkg/ingest"
	"github.com/shashank-mugiwara/laughingtale/pkg/type_configs"
	"github.com/shashank-mugiwara/laughingtale/pkg/utils"
)

func InitMasterPoller() {
	// First time load all data from given source config
	GetAllLoaderSourceConfigs()
}

func GetAllLoaderSourceConfigs() {
	var loaderScenarioConfigs []type_configs.SourceConfigs
	result := db.GetlaughingtaleDb().Find(&loaderScenarioConfigs)
	if result.Error != nil {
		logger.GetLaughingTaleLogger().Error("Unable to get SourceConfigs from database. The Error is: ", result.Error.Error())
		return
	}

	if result.RowsAffected == 0 {
		logger.GetLaughingTaleLogger().Info("Number of SourceConfigs found: 0. Not running the master poller")
		return
	} else {
		logger.GetLaughingTaleLogger().Info("Number of SourceConfigs found: ", result.RowsAffected)
	}

	for _, scfg := range loaderScenarioConfigs {
		identifier := scfg.Identifier
		for _, cfg := range scfg.SourceConfig {
			go ProcessEachSourceConfig(identifier, cfg)
		}
	}
}

func ProcessEachSourceConfig(identifier string, sourceConfig type_configs.SourceConfig) {

	ctx := context.Background()
	entryKey := utils.GetSourceConfigStringRepresentation(identifier, sourceConfig)
	isInitialLoadReady, keyExistsErr := client.GetRedisClient().Get(ctx, entryKey).Result()

	if keyExistsErr != nil {
		logger.GetLaughingTaleLogger().Info(keyExistsErr)
	}

	if !utils.IsBlank(isInitialLoadReady) {
		return
	}

	// First time its always SIMPLE strategy based polling
	poller, err := factory.GetStrategyFactory("SIMPLE")
	if err != nil {
		logger.GetLaughingTaleLogger().Error(err.Error())
	}

	resultList, resultErr := poller.Poll(identifier, sourceConfig)
	if resultErr != nil {
		logger.GetLaughingTaleLogger().Error(resultErr.Error())
	}

	ingest.IngestPostgresToMongo(identifier, sourceConfig, resultList)

	redisErr := client.GetRedisClient().Set(ctx, entryKey, time.Now().UTC().String(), 0).Err()
	if redisErr != nil {
		logger.GetLaughingTaleLogger().Error("Unable to set entryKey to redis. Please check if redis instance is healty.")
	}

	lastUpdatedAtKey := utils.GetDatabaseNameWithCollectionName(identifier, sourceConfig)
	redisErr = client.GetRedisClient().Set(ctx, lastUpdatedAtKey, time.Now().UTC().String(), 0).Err()
	if redisErr != nil {
		logger.GetLaughingTaleLogger().Error("Unable to set lastUpdatedAtKey to redis. Please check if redis instance is healty.")
	}

}
