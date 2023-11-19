package poller_strategy

import (
	"fmt"
	"time"

	"github.com/shashank-mugiwara/laughingtale/db"
	"github.com/shashank-mugiwara/laughingtale/logger"
	"github.com/shashank-mugiwara/laughingtale/pkg/type_configs"
	"github.com/shashank-mugiwara/laughingtale/pkg/utils"
)

type SimpleStrategy struct {
	PollerStrategy
}

func newSimpleStrategyPoller() IPollerStrategy {
	return &SimpleStrategy{
		PollerStrategy: PollerStrategy{
			WhereQueryPrefix:         "",
			PollerFrequencyInSeconds: 60,
		},
	}
}

func (simpleStrategy *SimpleStrategy) Poll(identifier string, sourceConfig type_configs.SourceConfig) ([]interface{}, error) {
	resultList, resultErr := PollDataFromSource(identifier, sourceConfig)
	if resultErr != nil {
		return nil, resultErr
	}
	return resultList, resultErr
}

func PollDataFromSource(identifier string, sourceConfig type_configs.SourceConfig) ([]interface{}, error) {
	listOfFields, err := utils.GenerateListOfFields(sourceConfig)
	if err != nil {
		logger.GetLaughingTaleLogger().Error("Failed to poll data from database.")
		return nil, err
	}

	var query string

	if !utils.IsBlank(sourceConfig.FilterConfig.WhereQuery) {
		query = "select " + sourceConfig.PrimaryKey + " as primaryKeyId, " + listOfFields + " FROM " + sourceConfig.DbSchema + "." + sourceConfig.TableName + " WHERE " +
			sourceConfig.FilterConfig.WhereQuery + " LIMIT " + sourceConfig.FilterConfig.Limit + ";"
	} else {
		query = "select " + sourceConfig.PrimaryKey + " as primaryKeyId, " + listOfFields + " FROM " + sourceConfig.DbSchema + "." + sourceConfig.TableName + " LIMIT " + sourceConfig.FilterConfig.Limit + ";"
	}

	logger.GetLaughingTaleLogger().Info("Executing query: ", query)

	start := time.Now().UnixNano() / int64(time.Millisecond)

	rows, err := db.GetlaughingtaleDb().Raw(query).Rows()
	defer rows.Close()

	end := time.Now().UnixNano() / int64(time.Millisecond)
	diff := end - start
	fmt.Printf("Duration for loading data from postgres in iteration-1 (ms): %d ms\n", diff)
	resultList, resultErr := utils.ParseSqlRows(rows)
	if resultErr != nil {
		logger.GetLaughingTaleLogger().Error("Unable to poll and ingest data from given scenarioConfig. Error is: ", resultErr)
		return nil, resultErr
	}

	return resultList, resultErr
}
