package utils

import (
	"github.com/shashank-mugiwara/laughingtale/db"
	"github.com/shashank-mugiwara/laughingtale/logger"
	"github.com/shashank-mugiwara/laughingtale/pkg/type_configs"
)

type RecordCountResult struct {
	RecCount int `json:"rec_count"`
}

func ListTables(schema string) []string {
	var tables []string
	if err := db.GetlaughingtaleDb().Table("information_schema.tables").Where("table_schema = ?", "public").Pluck("table_name", &tables).Error; err != nil {
		panic(err)
	}
	return tables
}

func GetRecordCount(sourceConfig type_configs.SourceConfig) int {
	query := "SELECT COUNT(" + sourceConfig.PrimaryKey + ") as rec_count FROM " + sourceConfig.DbSchema + "." + sourceConfig.TableName + " WHERE " + sourceConfig.FilterConfig.WhereQuery
	recordNumberResult := RecordCountResult{RecCount: 0}
	db.GetlaughingtaleDb().Raw(query).Scan(&recordNumberResult)
	logger.GetLaughingTaleLogger().Info("Found ", recordNumberResult.RecCount, " number of records from table: ", sourceConfig.TableName)
	return recordNumberResult.RecCount
}
