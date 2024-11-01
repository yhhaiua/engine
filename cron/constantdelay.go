package cron

import "time"

// TimingSchedule represents a simple recurring duty cycle, e.g. "Every 5 minutes".
// It does not support jobs more frequent than once a second.
type TimingSchedule struct {
	Delay time.Duration
}

// Every returns a crontab Schedule that activates once every duration.
// Delays of less than a second are not supported (will round up to 1 second).
// Any fields less than a Second are truncated.
func Every(duration time.Duration) *TimingSchedule {
	if duration < time.Second {
		duration = time.Second
	}
	return &TimingSchedule{
		Delay: duration,
	}
}

// Next returns the next time this should be run.
// This rounds so that the next activation time will be on the second.
func (schedule *TimingSchedule) Next(t time.Time) time.Time {
	return t.Add(schedule.Delay)
}

func (schedule *TimingSchedule) Once() bool {
	return false
}
