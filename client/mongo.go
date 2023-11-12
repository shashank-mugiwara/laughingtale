package client

import (
	"context"
	"log"
	"time"

	"github.com/shashank-mugiwara/laughingtale/conf"
	"github.com/shashank-mugiwara/laughingtale/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var mongoDb *mongo.Client

func GetMongoClient() *mongo.Client {
	return mongoDb
}

func InitMongodb() {

	if conf.MongoSetting.Enabled == "false" {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().
		ApplyURI("mongodb://"+conf.MongoSetting.MasterNode+":"+conf.MongoSetting.Port))

	if err != nil {
		log.Fatalln("Unable to connect to MongoDB node.")
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalln("Unable to ping MongoDB, make sure it is running.")
	}

	logger.GetLaughingTaleLogger().Info("Connection to MongoDb is successful.")
	mongoDb = client
}
