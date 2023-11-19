package utils

import (
	"strings"

	"github.com/shashank-mugiwara/laughingtale/logger"
	"github.com/shashank-mugiwara/laughingtale/pkg/type_configs"
)

func GetSourceConfigStringRepresentation(identifier string, sourceConfig type_configs.SourceConfig) string {
	return identifier + ":" + sourceConfig.DbSchema + ":" + sourceConfig.TableName + ":" + sourceConfig.PrimaryKey + ":" + identifier + ":" + sourceConfig.TargetCollectionName + ":" +
		strings.Join(sourceConfig.ColumnList, ",") + ":" + sourceConfig.PollerConfig.PollingStrategy + ":" + sourceConfig.FilterConfig.WhereQuery + ":" + sourceConfig.FilterConfig.Limit + sourceConfig.Version
}

func GetDatabaseNameWithCollectionName(identifier string, sourceConfig type_configs.SourceConfig) string {
	return identifier + ":" + sourceConfig.TableName + ":" + sourceConfig.Version
}

func GenerateListOfFields(sourceConfig type_configs.SourceConfig) (string, error) {

	if len(sourceConfig.ColumnList) < 1 {
		logger.GetLaughingTaleLogger().Info("No columns specified, selecting all columns by default")
		return "*", nil
	}

	var selectedFields string = ""
	for _, field := range sourceConfig.ColumnList {
		selectedFields = selectedFields + "," + field
	}

	if len([]rune(selectedFields)) < 2 {
		logger.GetLaughingTaleLogger().Error("Failed to construct query string, so selecting all columns by default")
		return "*", nil
	}

	return selectedFields[1:], nil
}
