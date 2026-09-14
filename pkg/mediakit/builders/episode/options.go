package episode

import "time"

type EndChapterMode string

const (
	EndChapterModeClose EndChapterMode = "close"
	EndChapterModePeek  EndChapterMode = "peek"
)

type ShortChapterMode string

const (
	ShortChapterModeNone    ShortChapterMode = "none"
	ShortChapterModeDiscard ShortChapterMode = "discard"
	ShortChapterModeInclude ShortChapterMode = "include"
)

type Options struct {
	// Length of chapter to consider as end of episode
	// (chapters longer than this will continue the episode)
	EndingChapterDuration time.Duration

	// Sets ending chapter behavior
	EndingChapterMode EndChapterMode

	// Skips check on EndingChapterDuration
	IgnoreMissingEnd bool

	// Minimum number of chapters to constitute an episode
	MinimumChapters int

	// Minimum duration to be considered a valid episode
	// -- episodes shorter than this will be discarded
	// (such as the ending bits for DVD credits, etc.)
	MinimumEpisodeDuration time.Duration

	// Short chapter duration defines the maximum length of a "short" chapter.
	// These might be next episode previews, or similar.
	ShortChapterDuration time.Duration

	// ShortChapterMode defines how to handle a short chapter.
	// - none: no special handling (default)
	// - discard: ends episode immediately without adding, then skips
	//            to next non-short chapter
	// - include: adds the short chapter into the episode without
	//            end-of-episode processing
	ShortChapterMode ShortChapterMode
}

type OptionsBuilder struct {
	Options *Options
}

func NewOptionsBuilder() *OptionsBuilder {
	return &OptionsBuilder{
		Options: &Options{},
	}
}

func (b *OptionsBuilder) Build() Options {
	opts := new(Options)
	*opts = *b.Options

	if opts.EndingChapterMode == "" {
		opts.EndingChapterMode = EndChapterModeClose
	}

	if opts.ShortChapterMode == "" {
		opts.ShortChapterMode = ShortChapterModeNone
	}

	return *b.Options
}

func (b *OptionsBuilder) Copy() *OptionsBuilder {
	nb := NewOptionsBuilder()
	*nb.Options = *b.Options
	return nb
}

func (b *OptionsBuilder) EndingChapterDuration(d time.Duration) *OptionsBuilder {
	nb := b.Copy()
	nb.Options.EndingChapterDuration = d
	return nb
}

func (b *OptionsBuilder) EndingChapterMode(mode EndChapterMode) *OptionsBuilder {
	nb := b.Copy()
	nb.Options.EndingChapterMode = mode
	return nb
}

func (b *OptionsBuilder) IgnoreMissingEnd(e bool) *OptionsBuilder {
	nb := b.Copy()
	nb.Options.IgnoreMissingEnd = e
	return nb
}

func (b *OptionsBuilder) MinimumChapters(s int) *OptionsBuilder {
	nb := b.Copy()
	nb.Options.MinimumChapters = s
	return nb
}

func (b *OptionsBuilder) MinimumEpisodeDuration(d time.Duration) *OptionsBuilder {
	nb := b.Copy()
	nb.Options.MinimumEpisodeDuration = d
	return nb
}

func (b *OptionsBuilder) ShortChapterDuration(d time.Duration) *OptionsBuilder {
	nb := b.Copy()
	nb.Options.ShortChapterDuration = d
	return nb
}

func (b *OptionsBuilder) ShortChapterMode(m ShortChapterMode) *OptionsBuilder {
	nb := b.Copy()
	nb.Options.ShortChapterMode = m
	return nb
}
