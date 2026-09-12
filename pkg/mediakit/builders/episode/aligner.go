package episode

import (
	"io"
	"time"

	"github.com/go-logr/logr"

	"github.com/neptune-media/MediaKit-go/pkg/helpers"
	"github.com/neptune-media/MediaKit-go/pkg/mediakit"
	"github.com/neptune-media/MediaKit-go/pkg/mediakit/utils/timing"
)

// Aligner is used to modify the provided list of chapters to align
// them with I-frames to reduce video corruption when splitting.
type Aligner struct {
	Logger logr.Logger
}

func (a *Aligner) AlignChaptersToIFrames(chapters mediakit.ChapterList, seeker FrameSeeker) (mediakit.ChapterList, error) {
	defer timing.NewFuncDurationPrinter(time.Now(), "finished aligning chapters to I-frames")(a.Logger)

	aligned := make(mediakit.ChapterList, 0, len(chapters))

	for idx, chapter := range chapters {
		// Align the current chapter to a frame
		curr, err := a.alignChapter(chapter, seeker)

		// Exit if we hit an error
		if err := helpers.IgnoreEOF(err); err != nil {
			return nil, err
		}

		// If this isn't the first chapter, we need to align the previous chapter end
		// to our new start point.
		if idx > 0 {
			prev := aligned[idx-1]

			// Check if the previous chapter ends after this chapter starts
			// and adjust the previous end time to the new start time
			if prev.TimeEnd > curr.TimeStart {
				a.Logger.V(1).Info("updated chapter end time",
					"chapter", prev.ID,
					"old-end-time", prev.EndTime().Seconds(),
					"new-end-time", curr.StartTime().Seconds(),
				)
				prev.TimeEnd = curr.TimeStart

				// Store updated chapter
				aligned[idx-1] = prev
			}
		}

		// Store updated chapter
		aligned = append(aligned, curr)
	}

	return aligned, nil
}

func (a *Aligner) alignChapter(chapter mediakit.Chapter, seeker FrameSeeker) (mediakit.Chapter, error) {
	// Seek until the next frame is after the chapter start time
	for seeker.Current(); seeker.Current() < chapter.StartTime(); seeker.Next() {
		// Stop seeking if we're at the end
		if seeker.EOF() {
			return mediakit.Chapter{}, io.EOF
		}
	}

	// Create new chapter with updated time information
	aligned := chapter
	aligned.TimeStart = int64(seeker.Current() / time.Millisecond)
	a.Logger.V(1).Info("aligned chapter start",
		"chapter", chapter.ID,
		"old-start-time", chapter.StartTime().Seconds(),
		"new-start-time", aligned.StartTime().Seconds(),
	)

	return aligned, nil
}
