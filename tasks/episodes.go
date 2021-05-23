package tasks

import (
	"github.com/neptune-media/MediaKit-go"
	"time"
)

func ReadVideoEpisodes(filename string, opts mediakit.EpisodeBuilderOptions) ([]mediakit.Episode, error) {
	episodes := make([]mediakit.Episode, 0)
	chapters, err := ReadVideoChapters(filename)
	if err != nil {
		return nil, err
	}

	if opts.FrameSeeker != nil {
		AlignChaptersToIFrames(chapters, opts.FrameSeeker)
	}

	episode := mediakit.Episode{}
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
		episode = mediakit.Episode{}
	}

	return episodes, nil
}

func AlignChaptersToIFrames(chapters []mediakit.Chapter, seeker *mediakit.FrameSeeker) {
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
