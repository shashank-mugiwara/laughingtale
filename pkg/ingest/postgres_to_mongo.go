package ingest

import (
	"context"
	"fmt"
	"time"

	"github.com/shashank-mugiwara/laughingtale/client"
	"github.com/shashank-mugiwara/laughingtale/logger"
	"github.com/shashank-mugiwara/laughingtale/pkg/type_configs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func IngestPostgresToMongo(identifier string, sourceConfig type_configs.SourceConfig, resultList []interface{}) {
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
