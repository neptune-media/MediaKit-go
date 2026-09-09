package episode

import "time"

type Options struct {
	// Length of chapter to consider as end of episode
	// (chapters longer than this will continue the episode)
	EndingChapterDuration time.Duration

	// Skips check on EndingChapterDuration
	IgnoreMissingEnd bool

	// Minimum number of chapters to constitute an episode
	MinimumChapters int

	// Minimum duration to be considered a valid episode
	// -- episodes shorter than this will be discarded
	// (such as the ending bits for DVD credits, etc.)
	MinimumEpisodeDuration time.Duration
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
	b.Options.MinimumEpisodeDuration = d
	return nb
}
