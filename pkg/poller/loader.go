package poller

import (
	"time"

	"github.com/shashank-mugiwara/laughingtale/db"
	"github.com/shashank-mugiwara/laughingtale/logger"
	"github.com/shashank-mugiwara/laughingtale/pkg/ingest"
	"github.com/shashank-mugiwara/laughingtale/pkg/type_configs"
)

func InitMasterPoller() {
	go func() {
		for true {
			GetAllLoaderSourceConfigs()
			time.Sleep(1 * time.Minute)
		}
	}()
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
		PollData(scfg)
	}
}

func PollData(loaderSourceConfigs type_configs.SourceConfigs) {
	identifier := loaderSourceConfigs.Identifier
	for _, scfg := range loaderSourceConfigs.SourceConfig {
		go ingest.PollDataFromSourceAndIngestToTarget(identifier, scfg)
	}
}
