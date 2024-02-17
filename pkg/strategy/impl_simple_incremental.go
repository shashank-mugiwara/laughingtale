package poller_strategy

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/shashank-mugiwara/laughingtale/client"
	"github.com/shashank-mugiwara/laughingtale/db"
	"github.com/shashank-mugiwara/laughingtale/logger"
	"github.com/shashank-mugiwara/laughingtale/pkg/ingest"
	"github.com/shashank-mugiwara/laughingtale/pkg/type_configs"
	"github.com/shashank-mugiwara/laughingtale/pkg/utils"
)

type SimpleIncrementalStrategy struct {
	PollerStrategy
}

func newSimpleIncrementalStrategyPoller() IPollerStrategy {
	return &SimpleIncrementalStrategy{
		PollerStrategy: PollerStrategy{
			WhereQueryPrefix:         "updated_at > NOW() - INTERVAL '{incrementalPollFreq} minutes'",
			PollerFrequencyInSeconds: 60,
		},
	}
}

func (simpleIncremental *SimpleIncrementalStrategy) Poll(identifier string, sourceConfig type_configs.SourceConfig) ([]interface{}, error) {
	ctx := context.Background()
	lastUpdatedKey := utils.GetDatabaseNameWithCollectionName(identifier, sourceConfig)
	lastIncrementalUpdateTimestamp := client.GetRedisClient().Get(ctx, lastUpdatedKey)

	if lastIncrementalUpdateTimestamp == nil {
		logger.GetLaughingTaleLogger().Info("Unable to find lastIncrementalUpdatedAt keys from redis. Running the poller with SIMPLE strategy again for the given source config: ", utils.GetSourceConfigStringRepresentation(identifier, sourceConfig))
		simplePoll := &SimpleStrategy{}
		resultList, resultErr := simplePoll.Poll(identifier, sourceConfig)
		if resultErr != nil {
			ingest.IngestPostgresToMongo(identifier, sourceConfig, resultList)
		} else {
			return nil, errors.New("Unable to find lastIncrementalUpdatedAt keys from redis even after running master poller. Returning")
		}

		logger.GetLaughingTaleLogger().Info("Ran master poller successfully")
	}

	lastIncrementalUpdateTimestamp = client.GetRedisClient().Get(ctx, lastUpdatedKey)
	if lastIncrementalUpdateTimestamp == nil {
		logger.GetLaughingTaleLogger().Error("Unable to find lastIncrementalUpdatedAt keys from redis even after running master poller. Returning")
		return nil, errors.New("Unable to find lastIncrementalUpdatedAt keys from redis even after running master poller. Returning")
	}

	listOfFields, err := utils.GenerateListOfFields(sourceConfig)
	if err != nil {
		logger.GetLaughingTaleLogger().Error("Failed to poll data from database.")
		return nil, err
	}

	var query string

	if !utils.IsBlank(sourceConfig.FilterConfig.WhereQuery) {
		query = "select " + sourceConfig.PrimaryKey + " as primaryKeyId, " + listOfFields + " FROM " + sourceConfig.DbSchema + "." + sourceConfig.TableName + " WHERE " + strings.Replace(simpleIncremental.WhereQueryPrefix, "{incrementalPollFreq}", sourceConfig.PollerConfig.DeltaUpdateIntervalInMinutes, 1)
	} else {
		query = "select " + sourceConfig.PrimaryKey + " as primaryKeyId, " + listOfFields + " FROM " + sourceConfig.DbSchema + "." + sourceConfig.TableName + " WHERE " + strings.Replace(simpleIncremental.WhereQueryPrefix, "{incrementalPollFreq}", sourceConfig.PollerConfig.DeltaUpdateIntervalInMinutes, 1)
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
