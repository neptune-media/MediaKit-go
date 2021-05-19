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

func ReadVideoEpisodes(filename string, opts EpisodeBuilderOptions) ([]Episode, error) {
	episodes := make([]Episode, 0)
	chapters, err := ReadVideoChapters(filename)
	if err != nil {
		return nil, err
	}

	if opts.FrameSeeker != nil {
		AlignChaptersToIFrames(chapters, opts.FrameSeeker)
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

		// Discard the episode if the runtime doesn't meet the minimum length
		if episode.Runtime() < opts.MinimumEpisodeLength {
			episode.Discard = true
		}

		episodes = append(episodes, episode)
		episode = Episode{}
	}

	return episodes, nil
}

func AlignChaptersToIFrames(chapters []Chapter, seeker *FrameSeeker) {
	for chIndex, chapter := range chapters {
		for {
			// Stop seeking if at the end
			if seeker.AtEnd() {
				break
			}

			// Seek until the next frame is after the chapter start time
			if seeker.Peek() < chapter.StartTime() {
				seeker.Next()
				continue
			}

			// Update start time to current frame
			newStartTime := seeker.Current()
			//fmt.Printf("Updating chapter %d start time from %.1f to %.1f\n", chIndex, chapter.StartTime().Seconds(), newStartTime.Seconds())
			chapter.TimeStart = int64(newStartTime / time.Millisecond)

			// If this is not the first chapter
			if chIndex > 0 {
				prevChapter := chapters[chIndex-1]

				// Check if the previous chapter ends after this chapter starts
				// and adjust the previous end time to the new start time
				if prevChapter.TimeEnd > chapter.TimeStart {
					//fmt.Printf("Updating chapter %d end time from %.1f to %.1f\n", chIndex-1, prevChapter.EndTime().Seconds(), newStartTime.Seconds())
					prevChapter.TimeEnd = chapter.TimeStart
				}
			}
			break
		}
	}
}

func (e Episode) Runtime() time.Duration {
	r := time.Duration(0)
	for _, c := range e.Chapters {
		r += c.Runtime()
	}

	return r
}
