package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"

	"github.com/shashank-mugiwara/laughingtale/conf"
	"github.com/shashank-mugiwara/laughingtale/initialize"
	"github.com/shashank-mugiwara/laughingtale/pkg/poller"
	"github.com/shashank-mugiwara/laughingtale/pkg/type_config"
)

func main() {
	initialize.InitConfig()
	conf.InitEngine()
	app := conf.GetLaughingTaleEngine()
	initialize.InitRoutes()
	initialize.InitClients()

	sourceConfigContainer := &type_config.SourceConfigContainer{
		Identifier: "randomIdentifier",
		SourceConfig: []type_config.SourceConfig{
			{
				TargetDatabaseName:   "shield",
				TargetCollectionName: "app_form",
				DbSchema:             "shield",
				TableName:            "app_form",
				PrimaryKey:           "id",
				PrimaryKeyType:       "string",
				ColumnList:           []string{},
				FilterConfig: type_config.FilterConfig{
					WhereQuery: "created_at >= NOW() - INTERVAL '400 days'",
					Limit:      "100000",
				},
			},
			{
				TargetDatabaseName:   "shield",
				TargetCollectionName: "applicant",
				DbSchema:             "shield",
				TableName:            "applicant",
				PrimaryKey:           "id",
				PrimaryKeyType:       "int64",
				ColumnList:           []string{},
				FilterConfig: type_config.FilterConfig{
					WhereQuery: "created_at >= NOW() - INTERVAL '400 days'",
					Limit:      "100000",
				},
			},
			{
				TargetDatabaseName:   "groot",
				TargetCollectionName: "loan_product",
				DbSchema:             "groot",
				TableName:            "loan_product",
				PrimaryKey:           "id",
				PrimaryKeyType:       "int64",
				ColumnList:           []string{},
				FilterConfig: type_config.FilterConfig{
					Limit: "100000",
				},
			},
			{
				TargetDatabaseName:   "heimdall",
				TargetCollectionName: "roles",
				DbSchema:             "heimdall",
				TableName:            "roles",
				PrimaryKey:           "id",
				PrimaryKeyType:       "int64",
				ColumnList:           []string{},
				FilterConfig: type_config.FilterConfig{
					Limit: "100000",
				},
			},
		},
	}

	poller.PollData(sourceConfigContainer)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		fmt.Println("Gracefully shutting down...")
		_ = app.Shutdown()
	}()

	if err := app.Listen("0.0.0.0:" + strconv.Itoa(conf.ServerSetting.HttpPort)); err != nil {
		log.Panic(err)
	}

	fmt.Println("Running cleanup tasks...")
}
