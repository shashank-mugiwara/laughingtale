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
)

func main() {
	initialize.InitConfig()
	conf.InitEngine()
	app := conf.GetLaughingTaleEngine()
	initialize.InitClients()
	initialize.InitRoutes()
	poller.InitMasterPoller()

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
