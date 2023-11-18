package poller

import (
	"github.com/shashank-mugiwara/laughingtale/db"
	"github.com/shashank-mugiwara/laughingtale/logger"
	"github.com/shashank-mugiwara/laughingtale/pkg/type_configs"
)

func InitMasterPoller() {

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

}

func PollData(sourceConfig *type_configs.SourceConfigsDto) {
	var lstOfSources []string
	for _, sc := range sourceConfig.SourceConfig {
		lstOfSources = append(lstOfSources, sc.DbSchema+"."+sc.TableName)
	}

	logger.GetLaughingTaleLogger().Info("List of sources: ", lstOfSources)
}
