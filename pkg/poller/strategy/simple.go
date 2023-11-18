package poller_strategy

import "database/sql"

type SimpleStrategy struct {
	PollerStrategy
}

func newSimpleStrategyPoller() IPollerStrategy {
	return &SimpleStrategy{
		PollerStrategy: PollerStrategy{
			WhereQueryPrefix:         "",
			PollerFrequencyInSeconds: 60,
		},
	}
}

func (simpleStrategy *SimpleStrategy) Poll() *sql.Rows {
	return nil
}
