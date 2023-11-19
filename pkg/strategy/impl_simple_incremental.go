package poller_strategy

import "github.com/shashank-mugiwara/laughingtale/pkg/type_configs"

type SimpleIncrementalStrategy struct {
	PollerStrategy
}

func newSimpleIncrementalStrategyPoller() IPollerStrategy {
	return &SimpleIncrementalStrategy{
		PollerStrategy: PollerStrategy{
			WhereQueryPrefix:         "",
			PollerFrequencyInSeconds: 60,
		},
	}
}

func (simpleIncremental *SimpleIncrementalStrategy) Poll(identifier string, sourceConfig type_configs.SourceConfig) ([]interface{}, error) {
	return nil, nil
}
