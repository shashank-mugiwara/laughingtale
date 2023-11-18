package poller_strategy

import "database/sql"

type IPollerStrategy interface {
	Poll() *sql.Rows
}
