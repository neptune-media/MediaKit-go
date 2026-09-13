package mediakit

import "time"

// Episode is a collection of chapters in a multi-episode video file
type Episode struct {
	// Chapters is a list of chapters that make an episode
	Chapters ChapterList

	// Discard - when true, don't write episode to a file
	Discard bool
}

type EpisodeList []Episode

func (e Episode) EndTime() time.Duration {
	if len(e.Chapters) == 0 {
		return 0
	}

	return e.Chapters[len(e.Chapters)-1].EndTime()
}

func (e Episode) Runtime() time.Duration {
	var r time.Duration
	for _, c := range e.Chapters {
		r += c.Runtime()
	}

	return r
}

func (e Episode) StartTime() time.Duration {
	if len(e.Chapters) == 0 {
		return 0
	}

	return e.Chapters[0].StartTime()
}

func (l EpisodeList) Runtime() time.Duration {
	var runtime time.Duration
	if l == nil {
		return 0
	}

	for _, episode := range l {
		runtime += episode.Runtime()
	}

	return runtime
}
