package poller_strategy

import (
	"github.com/shashank-mugiwara/laughingtale/pkg/type_configs"
)

type IPollerStrategy interface {
	Poll(identifier string, sourceConfig type_configs.SourceConfig) ([]interface{}, error)
}

type PollerStrategy struct {
	WhereQueryPrefix         string
	PollerFrequencyInSeconds int
}
