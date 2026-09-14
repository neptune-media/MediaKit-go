package cmd

import (
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/neptune-media/MediaKit-go/cmd/common"
)

// splitCmd represents the split command
var splitCmd = &cobra.Command{
	Use:   "split [file]",
	Short: "Splits a multi-episode file into multiple files",
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

		// Split episodes
		err = common.SplitEpisodes(cmd.Context(), logger, filename, episodes)
		if err != nil {
			logger.Fatal("failed to split episodes", zap.Error(err))
		}

		// Fix chapter names
		err = common.FixChapterNames(cmd.Context(), logger, "split.mkv", episodes)
		if err != nil {
			logger.Fatal("failed to fix chapter names", zap.Error(err))
		}
	},
}

func init() {
	rootCmd.AddCommand(splitCmd)

	splitCmd.Flags().Bool(common.ArgAlignChapters, false, "Align chapters to I-Frames to reduce video corruption on split")
	splitCmd.Flags().String(common.ArgFramesFile, "", "path to frame parquet file")
	splitCmd.Flags().Duration(common.ArgEndingChapterDuration, time.Minute, "Chapters longer than this will continue the episode")
	splitCmd.Flags().Bool(common.ArgIgnoreMissingEnd, false, "Ignore missing end of episodes")
	splitCmd.Flags().Int(common.ArgMinChapters, 2, "Minimum number of chapters in an episode")
	splitCmd.Flags().Duration(common.ArgMinEpisodeDuration, 20*time.Minute, "Minimum runtime of an episode")
	splitCmd.Flags().Duration(common.ArgShortChapterDuration, 30*time.Second, "Defines max length of short chapters")
	splitCmd.Flags().String(common.ArgShortChapterMode, "none", "How to handle short chapters (none, discard, include)")
}
