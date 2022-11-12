package mediakit

import (
	"time"
)

// Episode is a collection of chapters in a multi-episode video file
type Episode struct {
	// Chapters is a list of chapters that make an episode
	Chapters []Chapter
	// Discard - when true, don't write episode to a file
	Discard bool
}

// EpisodeBuilderOptions is used to provide hints to the builder
// when processing a multi-episode video file.
type EpisodeBuilderOptions struct {
	// Last added chapter must exceed this to continue adding chapters
	EndingChapterTime time.Duration

	// Used to align chapters to frames, when provided
	FrameSeeker *FrameSeeker

	// Skips check on EndingChapterTime
	IgnoreMissingEnd bool

	// Minimum number of chapters to constitute an episode
	MinimumChapters int

	// Minimum length to be considered a valid episode
	// -- episodes shorter than this will be discarded
	// (such as the ending bits for DVD credits, etc)
	MinimumEpisodeLength time.Duration
}

func (e Episode) Runtime() time.Duration {
	r := time.Duration(0)
	for _, c := range e.Chapters {
		r += c.Runtime()
	}

	return r
}
