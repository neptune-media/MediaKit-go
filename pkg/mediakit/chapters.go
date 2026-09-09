package mediakit

import (
	"time"
)

// Chapter represents a chapter within a video file
type Chapter struct {
	// ID is the chapter id provided by the video file
	ID uint64

	// TimeStart is the start time of the chapter.  The units
	// are determined by the timescale provided by the
	// video.
	TimeStart int64

	// TimeEnd is the end time of the chapter.  The units
	// are determined by the timescale provided by the
	// video.
	TimeEnd int64

	// Enabled is a flag indicated if the chapter is enabled
	// in the video file.
	Enabled bool

	// Title is the encoded chapter title name
	Title string
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

func (c Chapter) WithTimescale(ts time.Duration) Chapter {
	t := c

	t.TimeEnd = t.TimeEnd / int64(ts)
	t.TimeStart = t.TimeStart / int64(ts)

	return t
}

type ChapterList []Chapter

func (l ChapterList) First() Chapter {
	if l == nil {
		return Chapter{}
	}

	return l[0]
}

func (l ChapterList) Last() Chapter {
	if l == nil {
		return Chapter{}
	}

	return l[len(l)-1]
}

func (l ChapterList) Runtime() time.Duration {
	var runtime time.Duration
	if l == nil {
		return 0
	}

	for _, ch := range l {
		runtime += ch.Runtime()
	}

	return runtime
}
