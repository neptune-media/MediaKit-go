package episode

import (
	"github.com/go-logr/logr"

	"github.com/neptune-media/MediaKit-go/pkg/mediakit"
)

const (
	ActionNone                   FilterAction = iota // No action, move to next filter
	ActionAppendChapter                              // Appends chapter to episode and moves to next chapter
	ActionCloseEpisode                               // Appends chapter and closes episode, then starts a new one
	ActionDiscardChapterAndClose                     // Discards the current chapter and closes the episode
	ActionDiscardEpisode                             // Discards the current episode and starts the next one
)

const (
	EpisodeSelectorCurrent EpisodeSelector = iota // Select current episode
	EpisodeSelectorNext                           // Select next episode
)

type FilterAction int

type FilterInput struct {
	Chapter       *mediakit.Chapter // Current chapter to consider
	ChapterIndex  int               // Index of the current chapter from the chapter list
	Episode       *mediakit.Episode // Current episode being built
	Logger        logr.Logger
	NextChapter   *mediakit.Chapter // Upcoming chapter, nil if the current is the last chapter
	Options       Options           // Options the builder was invoked with
	TotalChapters int               // Total number of chapters passed to the builder
}

// FirstChapter returns true if the current chapter is the first
// chapter of the ChapterList provided to the builder.
func (f FilterInput) FirstChapter() bool {
	return f.ChapterIndex == 0
}

// LastChapter returns true if the current chapter is the last
// chapter of the ChapterList provided to the builder.
func (f FilterInput) LastChapter() bool {
	return f.ChapterIndex+1 >= f.TotalChapters
}

type BuilderFilter func(input FilterInput) FilterAction

// EndingChapterFilter checks if the current chapter is the
// likely ending chapter of an episode.  Returns ActionCloseEpisode
// if true or ActionNone otherwise.
func EndingChapterFilter(input FilterInput) FilterAction {
	if input.Chapter.Runtime() <= input.Options.EndingChapterDuration {
		action := ActionCloseEpisode
		if input.Options.EndingChapterMode == EndChapterModePeek {
			if !input.LastChapter() {
				if input.NextChapter.Runtime() <= input.Options.ShortChapterDuration {
					action = ActionAppendChapter
				}
			}
		}
		return action
	}

	return ActionNone
}

// IgnoreMissingEndFilter checks if this is the last chapter,
// and closes the episode if true.  Does nothing otherwise.
func IgnoreMissingEndFilter(input FilterInput) FilterAction {
	// Last chapter, doesn't matter if we're the end or still missing one
	if input.LastChapter() {
		return ActionCloseEpisode
	}

	// Not the last chapter, we have no opinion
	return ActionNone
}

// MinimumChapterFilter checks if the current episode has at least
// the minimum number of chapters set in the builder options.
// Returns ActionAppendChapter if below the minimum, ActionNone otherwise.
func MinimumChapterFilter(input FilterInput) FilterAction {
	if len(input.Episode.Chapters) < input.Options.MinimumChapters {
		return ActionAppendChapter
	}

	return ActionNone
}

// MinimumEpisodeRuntimeFilter checks if the current episode runtime has
// reached the minimum set in the builder options.  Mutates Episode to
// set the Discard flag, returns ActionNone.
func MinimumEpisodeRuntimeFilter(input FilterInput) FilterAction {
	runtime := input.Episode.Runtime() + input.Chapter.Runtime()
	discard := runtime < input.Options.MinimumEpisodeDuration
	if discard != input.Episode.Discard {
		input.Logger.V(1).Info("mutating episode", "key", "discard", "old-value", input.Episode.Discard, "new-value", discard)
	}
	input.Episode.Discard = discard

	return ActionNone
}

type EpisodeSelector int

type ShortChapterIncludeOptions struct {
	EpisodeSelector EpisodeSelector
}

// ShortChapterFilter is used to handle short chapters.
func ShortChapterFilter(mode ShortChapterMode, includeOptions ShortChapterIncludeOptions) BuilderFilter {
	return func(input FilterInput) FilterAction {
		if input.Chapter.Runtime() > input.Options.ShortChapterDuration {
			return ActionNone
		}

		switch mode {
		case ShortChapterModeDiscard:
			// Discard mode is used to end the episode, drop the short
			// chapter, and start on the next.
			return ActionDiscardChapterAndClose
		case ShortChapterModeInclude:
			// Include mode is used to add the short chapter to the episode,
			// and keep building without ending the episode.
			if includeOptions.EpisodeSelector == EpisodeSelectorNext {
				return ActionCloseEpisode
			}
			return ActionAppendChapter
		default:
			return ActionNone
		}
	}
}
