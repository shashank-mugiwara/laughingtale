package ingest

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/shashank-mugiwara/laughingtale/client"
	"github.com/shashank-mugiwara/laughingtale/db"
	"github.com/shashank-mugiwara/laughingtale/logger"
	"github.com/shashank-mugiwara/laughingtale/pkg/type_configs"
	"github.com/shashank-mugiwara/laughingtale/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func PollDataFromSourceAndIngestToTarget(identifier string, sourceConfig type_configs.SourceConfig) {
	listOfFields, err := generateListOfFields(sourceConfig)
	if err != nil {
		logger.GetLaughingTaleLogger().Error("Failed to poll data from database.")
		return
	}

	utils.GetRecordCount(sourceConfig)

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
	resultList, resultErr := parse(rows)
	if resultErr != nil {
		logger.GetLaughingTaleLogger().Error("Unable to poll and ingest data from given scenarioConfig. Error is: ", resultErr)
		return
	}

	go ingestDataToTarget(identifier, sourceConfig, resultList)
}

func ingestDataToTarget(identifier string, sourceConfig type_configs.SourceConfig, resultList []interface{}) {
	mongoCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	// For now we are using on-container mongodb as target database

	targetDatabaseName := identifier
	database := client.GetMongoClient().Database(targetDatabaseName)
	collection := database.Collection(sourceConfig.TargetCollectionName)

	start := time.Now().UnixNano() / int64(time.Millisecond)
	var models []mongo.WriteModel

	for _, doc := range resultList {

		if sourceConfig.PrimaryKeyType == "string" {
			primaryKey := doc.(map[string]interface{})[sourceConfig.PrimaryKey].(string)
			model := []mongo.WriteModel{
				mongo.NewReplaceOneModel().SetFilter(bson.D{{Key: "_id", Value: primaryKey}}).
					SetReplacement(doc).SetUpsert(true),
			}
			models = append(models, model[0])
		} else if sourceConfig.PrimaryKeyType == "int64" {
			primaryKey := doc.(map[string]interface{})[sourceConfig.PrimaryKey].(int64)
			model := []mongo.WriteModel{
				mongo.NewReplaceOneModel().SetFilter(bson.D{{Key: "_id", Value: primaryKey}}).
					SetReplacement(doc).SetUpsert(true),
			}
			models = append(models, model[0])
		}

	}

	end := time.Now().UnixNano() / int64(time.Millisecond)
	diff := end - start
	fmt.Printf("Time for preparing the insert/updates for BulkWrites: %d ms\n", diff)

	opts := options.BulkWrite().SetOrdered(false)

	start = time.Now().UnixNano() / int64(time.Millisecond)
	_, insertErr := collection.BulkWrite(mongoCtx, models, opts)
	if insertErr != nil {
		logger.GetLaughingTaleLogger().Error("Erorr while bulk inseting data into on-container mongo: ", insertErr.Error())
	} else {
		logger.GetLaughingTaleLogger().Info("Successfully inserted documents to on-container mongodb")
	}
	end = time.Now().UnixNano() / int64(time.Millisecond)
	diff = end - start
	fmt.Printf("Time taken for ingesting the documents: %d ms\n", diff)
}

func generateListOfFields(sourceConfig type_configs.SourceConfig) (string, error) {

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

func parse(rows *sql.Rows) ([]interface{}, error) {
	columnTypes, err := rows.ColumnTypes()

	if err != nil {
		return nil, err
	}

	count := len(columnTypes)
	finalRows := []interface{}{}

	for rows.Next() {

		scanArgs := make([]interface{}, count)

		for i, v := range columnTypes {

			switch v.DatabaseTypeName() {
			case "VARCHAR", "TEXT", "UUID", "TIMESTAMP":
				scanArgs[i] = new(sql.NullString)
				break
			case "BOOL":
				scanArgs[i] = new(sql.NullBool)
				break
			case "INT4":
				scanArgs[i] = new(sql.NullInt64)
				break
			default:
				scanArgs[i] = new(sql.NullString)
			}
		}

		err := rows.Scan(scanArgs...)

		if err != nil {
			return nil, err
		}

		masterData := map[string]interface{}{}

		for i, v := range columnTypes {

			if z, ok := (scanArgs[i]).(*sql.NullBool); ok {
				masterData[v.Name()] = z.Bool
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullString); ok {
				masterData[v.Name()] = z.String
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullInt64); ok {
				masterData[v.Name()] = z.Int64
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullFloat64); ok {
				masterData[v.Name()] = z.Float64
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullInt32); ok {
				masterData[v.Name()] = z.Int32
				continue
			}

			masterData[v.Name()] = scanArgs[i]
		}

		finalRows = append(finalRows, masterData)
	}

	return finalRows, err
}
