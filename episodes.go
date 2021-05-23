package mediakit

import (
	"time"
)

type Episode struct {
	Chapters []Chapter
	Discard  bool
}

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
