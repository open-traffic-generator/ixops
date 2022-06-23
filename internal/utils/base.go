package utils

import "time"

type WaitForOpts struct {
	Condition string
	Interval  time.Duration
	Timeout   time.Duration
}
