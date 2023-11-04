package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"

	"github.com/shashank-mugiwara/laughingtale/conf"
	"github.com/shashank-mugiwara/laughingtale/initialize"
)

func main() {
	conf.InitEngine()
	app := conf.GetLaughingTaleEngine()
	initialize.InitRoutes()

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
