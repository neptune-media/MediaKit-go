package mediakit

import "time"

type Episode struct {
	Chapters []Chapter
	Discard  bool
}

type EpisodeBuilderOptions struct {
	// Last added chapter must exceed this to continue adding chapters
	EndingChapterTime time.Duration

	// Skips check on EndingChapterTime
	IgnoreMissingEnd bool

	// Minimum number of chapters to constitute an episode
	MinimumChapters int

	// Minimum length to be considered a valid episode
	// -- episodes shorter than this will be discarded
	// (such as the ending bits for DVD credits, etc)
	MinimumEpisodeLength time.Duration
}

func ReadVideoEpisodes(filename string, opts EpisodeBuilderOptions) ([]Episode, error) {
	episodes := make([]Episode, 0)
	chapters, err := ReadVideoChapters(filename)
	if err != nil {
		return nil, err
	}

	episode := Episode{}
	for _, chapter := range chapters {
		episode.Chapters = append(episode.Chapters, chapter)

		// Keep adding chapters till we meet the minimum
		if len(episode.Chapters) < opts.MinimumChapters {
			continue
		}

		// Keep adding chapters if the last chapter added is above the threshold
		if !opts.IgnoreMissingEnd && chapter.Runtime() > opts.EndingChapterTime {
			continue
		}

		if episode.Runtime() < opts.MinimumEpisodeLength {
			episode.Discard = true
		}

		episodes = append(episodes, episode)
		episode = Episode{}
	}

	return episodes, nil
}

func (e Episode) Runtime() time.Duration {
	r := time.Duration(0)
	for _, c := range e.Chapters {
		r += c.Runtime()
	}

	return r
}
