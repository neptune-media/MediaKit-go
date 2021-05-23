package mediakit

import (
	"time"
)

type Chapter struct {
	ID        uint64
	TimeStart int64
	TimeEnd   int64
	Enabled   bool
}

func (c Chapter) EndTime() time.Duration {
	return time.Duration(c.TimeEnd) * time.Millisecond
}

func (c Chapter) Runtime() time.Duration {
	return c.EndTime() - c.StartTime()
}

func (c Chapter) StartTime() time.Duration {
	return time.Duration(c.TimeStart) * time.Millisecond
}
