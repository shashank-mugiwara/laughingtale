package client

import (
	"github.com/redis/go-redis/v9"
	"github.com/shashank-mugiwara/laughingtale/conf"
	"github.com/shashank-mugiwara/laughingtale/logger"
)

var redisClient *redis.Client

func GetRedisClient() *redis.Client {
	return redisClient
}

func InitRedisClient() {
	redisClientDb := redis.NewClient(&redis.Options{
		Addr: conf.RedisSetting.Host + ":" + conf.RedisSetting.Port,
	})

	logger.GetLaughingTaleLogger().Info("Connection to redis successful")
	redisClient = redisClientDb
}
