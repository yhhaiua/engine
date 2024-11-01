package cron

import "time"

// DelaySchedule represents a simple recurring duty cycle, e.g. "Every 5 minutes".
// It does not support jobs more frequent than once a second.
type DelaySchedule struct {
	Delay time.Duration
}

// EveryDelay returns a crontab Schedule that activates once every duration.
// Delays of less than a second are not supported (will round up to 1 second).
// Any fields less than a Second are truncated.
func EveryDelay(duration time.Duration) *DelaySchedule {
	if duration < time.Second {
		duration = time.Second
	}
	return &DelaySchedule{
		Delay: duration,
	}
}

// Next returns the next time this should be run.
// This rounds so that the next activation time will be on the second.
func (schedule *DelaySchedule) Next(t time.Time) time.Time {
	return t.Add(schedule.Delay)
}

func (schedule *DelaySchedule) Once() bool {
	return true
}
