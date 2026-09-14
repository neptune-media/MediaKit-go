package info

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/neptune-media/MediaKit-go/cmd/common"
)

// episodesCmd represents the episodes command
var episodesCmd = &cobra.Command{
	Use:   "episodes [file]",
	Short: "Scans chapters in a file and combines them into episodes",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Setup logging
		logger := common.InitLogger()
		defer logger.Sync()

		// Get input
		filename := args[0]
		logger = logger.With(zap.String("job", filepath.Base(filename)))
		episodes, err := common.BuildEpisodes(logger, filename)
		if err != nil {
			logger.Fatal("failed to build episodes", zap.Error(err))
		}

		logger.Info("dumping episodes")
		for idx, episode := range episodes {
			chapters := make([]string, len(episode.Chapters))
			for i, c := range episode.Chapters {
				chapters[i] = fmt.Sprintf("%d (%.01fs)", c.ID, c.Runtime().Seconds())
			}
			fmt.Printf("%02d\t%08d - %08d\t%20s\tdiscard=%t\tCh: %s\n",
				idx+1,
				episode.Chapters.First().TimeStart,
				episode.Chapters.Last().TimeEnd,
				episode.Runtime(),
				episode.Discard,
				strings.Join(chapters, " , "))
		}
	},
}

func init() {
	infoCmd.AddCommand(episodesCmd)

	episodesCmd.Flags().Bool(common.ArgAlignChapters, false, "Align chapters to I-Frames to reduce video corruption on split")
	episodesCmd.Flags().String(common.ArgFramesFile, "", "path to frame parquet file")
	episodesCmd.Flags().Duration(common.ArgEndingChapterDuration, time.Minute, "Chapters longer than this will continue the episode")
	episodesCmd.Flags().String(common.ArgEndingChapterMode, "close", "What to do at the end of a chapter (close, peek)")
	episodesCmd.Flags().Bool(common.ArgIgnoreMissingEnd, false, "Ignore missing end of episodes")
	episodesCmd.Flags().Int(common.ArgMinChapters, 2, "Minimum number of chapters in an episode")
	episodesCmd.Flags().Duration(common.ArgMinEpisodeDuration, 20*time.Minute, "Minimum runtime of an episode")
	episodesCmd.Flags().Duration(common.ArgShortChapterDuration, 30*time.Second, "Defines max length of short chapters")
	episodesCmd.Flags().String(common.ArgShortChapterMode, "none", "How to handle short chapters (none, discard, include)")
}
