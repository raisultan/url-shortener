package analytics

import "time"

type ClickEvent struct {
	URLAlias  string
	Timestamp time.Time
	UserAgent string
	IP        string
	Referrer  string
	Latency   time.Duration
}
