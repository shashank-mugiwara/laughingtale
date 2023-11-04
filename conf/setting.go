package conf

import (
	"log"

	"gopkg.in/ini.v1"
)

type Application struct {
	RunType   string
	IsEnvProd bool
}

var ApplicationSetting = &Application{}

type Server struct {
	RunMode  string
	HttpPort int
}

var ServerSetting = &Server{}

var cfg *ini.File

func SetUp(path string) {
	var err error
	var tempCfg *ini.File

	if path != "" {
		tempCfg, err = ini.Load(path)
	} else {
		tempCfg, err = ini.Load("conf/config.ini")
	}

	cfg = tempCfg
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf/config.ini': %v", err)
	}

	mapTo("application", ApplicationSetting)
	mapTo("server", ServerSetting)
}

func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
