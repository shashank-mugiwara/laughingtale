package utils

import (
	"strings"

	"github.com/shashank-mugiwara/laughingtale/pkg/type_configs"
)

func GetSourceConfigStringRepresentation(identifier string, sourceConfig type_configs.SourceConfig) string {
	return identifier + ":" + sourceConfig.DbSchema + ":" + sourceConfig.TableName + ":" + sourceConfig.PrimaryKey + ":" + identifier + ":" + sourceConfig.TargetCollectionName + ":" +
		strings.Join(sourceConfig.ColumnList, ",") + ":" + sourceConfig.PollerConfig.PollingStrategy + ":" + sourceConfig.FilterConfig.WhereQuery + ":" + sourceConfig.FilterConfig.Limit
}
