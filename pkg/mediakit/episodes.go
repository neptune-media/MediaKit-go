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

func (e Episode) Runtime() time.Duration {
	var r time.Duration
	for _, c := range e.Chapters {
		r += c.Runtime()
	}

	return r
}
