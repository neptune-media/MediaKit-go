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
	}
	b.Options.EndingChapterMode = EndChapterModePeek

	b.Logger.Info("building episode", "start-chapter", start)
	for idx, chapter := range b.chapters[start:len(b.chapters)] {
		input := FilterInput{
			Chapter:       &chapter,
			ChapterIndex:  start + idx,
			Episode:       &episode,
			Logger:        b.Logger.WithName("filter"),
			Options:       b.Options,
			TotalChapters: len(b.chapters),
		}
		b.chapterPos++

		if b.chapterPos < len(b.chapters) {
			input.NextChapter = &b.chapters[b.chapterPos]
		}

		filters := []BuilderFilter{
			ShortChapterFilter(input.Options.ShortChapterMode, ShortChapterIncludeOptions{EpisodeSelector: EpisodeSelectorCurrent}),
			MinimumChapterFilter,
			// IgnoreMissingEndFilter,
			MinimumEpisodeRuntimeFilter,
			EndingChapterFilter,
		}

		switch action := b.applyFilters(filters, input); action {
		case ActionNone, ActionAppendChapter:
			episode.Chapters = append(episode.Chapters, chapter)
		case ActionCloseEpisode:
			episode.Chapters = append(episode.Chapters, chapter)
			return episode, nil
		case ActionDiscardChapterAndClose:
			return episode, nil
		case ActionDiscardEpisode:
			episode.Discard = true
			return episode, nil
		}
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

func (b *Builder) applyFilters(filters []BuilderFilter, input FilterInput) FilterAction {
	for _, filter := range filters {
		switch action := filter(input); action {
		case ActionNone:
			continue
		default:
			return action
		}
	}

	return ActionNone
}
