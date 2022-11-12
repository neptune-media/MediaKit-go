package mediakit

import (
	"time"
)

// Chapter represents a chapter within a video file
type Chapter struct {
	// ID is the chapter id provided by the video file
	ID uint64

	// TimeStart is the start time of the chapter.  The units
	// are determined by the time scale provided by the
	// video.
	TimeStart int64

	// TimeEnd is the end time of the chapter.  The units
	// are determined by the time scale provided by the
	// video.
	TimeEnd int64

	// Enabled is a flag indicated if the chapter is enabled
	// in the video file.
	Enabled bool
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
