package episode

import (
	"errors"
	"io"
	"time"

	"github.com/go-logr/logr"

	"github.com/neptune-media/MediaKit-go/pkg/mediakit"
	"github.com/neptune-media/MediaKit-go/pkg/mediakit/utils/timing"
)

type Builder struct {
	Logger  logr.Logger
	Options Options

	// List of chapters we're building from
	chapters   mediakit.ChapterList
	chapterPos int // Index of next chapter
}

type BuilderOption func(*Builder)

func WithLogger(logger logr.Logger) BuilderOption {
	return func(builder *Builder) {
		builder.Logger = logger
	}
}

func WithOptions(options Options) BuilderOption {
	return func(builder *Builder) {
		builder.Options = options
	}
}

func NewEpisodeBuilder(opts ...BuilderOption) *Builder {
	builder := &Builder{
		Logger: logr.Discard(),
	}

	for _, opt := range opts {
		opt(builder)
	}

	return builder
}

func (b *Builder) WithChapters(chapters mediakit.ChapterList) *Builder {
	// Copy builder
	nb := new(Builder)
	*nb = *b

	nb.chapters = chapters
	nb.chapterPos = 0
	return nb
}

// Build will build and return the next episode
func (b *Builder) Build() (mediakit.Episode, error) {
	defer timing.NewFuncDurationPrinter(time.Now(), "finished building episode")(b.Logger)

	start := b.chapterPos
	episode := mediakit.Episode{
		Chapters: make([]mediakit.Chapter, 0),
		Discard:  true,
	}

	b.Logger.Info("building episode", "start-chapter", start)
	for _, chapter := range b.chapters[start:len(b.chapters)] {
		b.chapterPos++

		// Check how we should handle short chapters
		if chapter.Runtime() < b.Options.ShortChapterDuration {
			switch b.Options.ShortChapterMode {
			case ShortChapterModeNone:
				// None is the default, normal handling.
			case ShortChapterModeDiscard:
				// Discard mode is used to end the episode, drop the short chapter, and
				// start on the next.

				// Return this episode, we're done now
				return episode, nil

			case ShortChapterModeInclude:
				// Include mode is used to add the short chapter to the episode, and
				// keep building without ending the episode.
				continue
			}
		}

		// Add next chapter to episode
		episode.Chapters = append(episode.Chapters, chapter)

		// Keep adding chapters until we meet the minimum
		if len(episode.Chapters) < b.Options.MinimumChapters {
			continue
		}

		// Keep adding chapters if the last chapter added is above the threshold
		if !b.Options.IgnoreMissingEnd && chapter.Runtime() > b.Options.EndingChapterDuration {
			continue
		}

		// // Check if we should handle short chapters
		// if b.Options.ShortChapterMode != ShortChapterModeNone && chapter.Runtime() < b.Options.ShortChapterDuration {
		// 	switch b.Options.ShortChapterMode {
		// 	case ShortChapterModeNone:
		// 		// None is the default, normal handling.
		// 	case ShortChapterModeDiscard:
		// 		// Discard mode is used to end the episode, drop the short chapter, and
		// 		// start on the next.
		// 		episode.Discard = false
		//
		// 		// Return this episode, we're done now
		// 		return episode, nil
		//
		// 	case ShortChapterModeInclude:
		// 		// Include mode is used to add the short chapter to the episode, and
		// 		// keep building without ending the episode.
		// 		continue
		// 	}
		// }

		// Discard the episode if the runtime doesn't meet the minimum length
		// Keep episode once it exceeds the minimum length
		if episode.Runtime() >= b.Options.MinimumEpisodeDuration {
			episode.Discard = false
		}

		// At this point, we're done building this episode
		return episode, nil
	}

	return episode, io.EOF
}

func (b *Builder) BuildAll() (mediakit.EpisodeList, error) {
	defer timing.NewFuncDurationPrinter(time.Now(), "finished building all episodes")(b.Logger)
	episodes := make(mediakit.EpisodeList, 0)
	building := true

	for building {
		episode, err := b.Build()
		if err != nil {
			building = false
			if !errors.Is(err, io.EOF) {
				return nil, err
			}
		}

		episodes = append(episodes, episode)
	}

	return episodes, nil
}
