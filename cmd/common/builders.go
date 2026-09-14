package common

import (
	"github.com/spf13/viper"

	"github.com/neptune-media/MediaKit-go/pkg/mediakit/builders/episode"
)

// NewEpisodeOptionsBuilder returns a new OptionsBuilder based
// on arguments provided on the command line.
func NewEpisodeOptionsBuilder() *episode.OptionsBuilder {
	return episode.NewOptionsBuilder().
		EndingChapterDuration(viper.GetDuration(ArgEndingChapterDuration)).
		EndingChapterMode(episode.EndChapterMode(viper.GetString(ArgEndingChapterMode))).
		IgnoreMissingEnd(viper.GetBool(ArgIgnoreMissingEnd)).
		MinimumChapters(viper.GetInt(ArgMinChapters)).
		MinimumEpisodeDuration(viper.GetDuration(ArgMinEpisodeDuration)).
		ShortChapterDuration(viper.GetDuration(ArgShortChapterDuration)).
		ShortChapterMode(episode.ShortChapterMode(viper.GetString(ArgShortChapterMode)))
}
