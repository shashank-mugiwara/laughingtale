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

type Kafka struct {
	Enabled string
	Host1   string
	Host2   string
	Port1   string
	Port2   string
}

type Mongo struct {
	Enabled    string
	MasterNode string
	Port       string
}

type Postgres struct {
	Enabled    string
	ReaderNode string
	Port       string
	Username   string
	Password   string
	Database   string
	SslEnabled string
}

var PostgresSetting = &Postgres{}

var ServerSetting = &Server{}

var KafkaSetting = &Kafka{}

var MongoSetting = &Mongo{}

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
	mapTo("kafka", KafkaSetting)
	mapTo("mongo", MongoSetting)
	mapTo("source-postgres", PostgresSetting)
}

func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
