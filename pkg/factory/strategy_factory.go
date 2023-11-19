package factory

import (
	"errors"

	poller_strategy "github.com/shashank-mugiwara/laughingtale/pkg/strategy"
)

func GetStrategyFactory(strategy_type string) (poller_strategy.IPollerStrategy, error) {
	if strategy_type == "SIMPLE" {
		return &poller_strategy.SimpleStrategy{}, nil
	}

	return nil, errors.New("No factory found for given factory type")
}
