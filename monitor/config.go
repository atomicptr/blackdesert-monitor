package monitor

import "time"

type Config struct {
	ProcessName  string
	PollInterval time.Duration
}
