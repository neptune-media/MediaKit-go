package info

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-logr/zapr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/neptune-media/MediaKit-go/cmd/common"
	"github.com/neptune-media/MediaKit-go/pkg/mediakit"
	"github.com/neptune-media/MediaKit-go/pkg/mediakit/builders/episode"
	"github.com/neptune-media/MediaKit-go/pkg/mediakit/utils"
	"github.com/neptune-media/MediaKit-go/pkg/tools/ffprobe"
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
		inputFilename := args[0]
		logger = logger.With(zap.String("job", filepath.Base(inputFilename)))
		logger.Info("using input file", zap.String("input-file", inputFilename))

		// Get frame data
		var frames ffprobe.FrameList = nil
		if viper.GetBool(common.ArgAlignChapters) {
			framesFilename := viper.GetString(common.ArgFramesFile)
			logger.Info("using frames file", zap.String("frames-file", framesFilename))

			frames, err = mediakit.ReadFramesFromFilename(framesFilename, zapr.NewLogger(logger), mediakit.FilterIFrames)
			if err != nil {
				logger.Fatal("failed to read frame data", zap.Error(err))
			}
		}

		// Read chapters
		logger.Info("reading chapters from file")
		f, err := os.Open(inputFilename)
		if err != nil {
			logger.Fatal("failed to open file", zap.Error(err))
		}
		defer f.Close()

		chapters, err := utils.ReadVideoChapters(f, zapr.NewLogger(logger))
		if err != nil {
			logger.Fatal("failed to read chapters from file", zap.Error(err))
		}

		if frames != nil {
			logger.Info("aligning chapters to I-frames")
			seeker := episode.NewFrameSeeker(frames)
			aligner := episode.Aligner{
				Logger: zapr.NewLogger(logger),
			}

			chapters, err = aligner.AlignChaptersToIFrames(chapters, seeker)
			if err != nil {
				logger.Fatal("failed to align chapters to I-frames", zap.Error(err))
			}
		}

		totalChapterRuntime := chapters.Runtime()

		logger.Info("building episodes from chapters")
		eob := episode.NewOptionsBuilder().
			EndingChapterDuration(viper.GetDuration(common.ArgEndingChapterDuration)).
			IgnoreMissingEnd(viper.GetBool(common.ArgIgnoreMissingEnd)).
			MinimumChapters(viper.GetInt(common.ArgMinChapters)).
			MinimumEpisodeDuration(viper.GetDuration(common.ArgMinEpisodeDuration)).
			ShortChapterDuration(viper.GetDuration(common.ArgShortChapterDuration)).
			ShortChapterMode(episode.ShortChapterMode(viper.GetString(common.ArgShortChapterMode)))

		builder := episode.NewEpisodeBuilder(
			episode.WithLogger(zapr.NewLogger(logger)),
			episode.WithOptions(eob.Build()),
		)
		episodes, err := builder.WithChapters(chapters).BuildAll()
		if err != nil {
			logger.Fatal("failed to build episodes", zap.Error(err))
		}

		logger.Info("dumping episodes")
		var totalEpisodeRuntime time.Duration
		for idx, episode := range episodes {
			totalEpisodeRuntime = totalEpisodeRuntime + episode.Runtime()
			fmt.Printf("%02d\t%08d - %08d\t%20s\tdiscard=%t\n",
				idx+1,
				episode.Chapters.First().TimeStart,
				episode.Chapters.Last().TimeEnd,
				episode.Runtime(),
				episode.Discard)
		}

		fmt.Println()
		fmt.Printf("pre-split run time: %20s\n", totalChapterRuntime)
		fmt.Printf("post-split run time: %20s\n", totalEpisodeRuntime)
		fmt.Println()
	},
}

func init() {
	infoCmd.AddCommand(episodesCmd)

	episodesCmd.Flags().Bool(common.ArgAlignChapters, false, "Align chapters to I-Frames to reduce video corruption on split")
	episodesCmd.Flags().String(common.ArgFramesFile, "", "path to frame parquet file")
	episodesCmd.Flags().Duration(common.ArgEndingChapterDuration, time.Minute, "Chapters longer than this will continue the episode")
	episodesCmd.Flags().Bool(common.ArgIgnoreMissingEnd, false, "Ignore missing end of episodes")
	episodesCmd.Flags().Int(common.ArgMinChapters, 2, "Minimum number of chapters in an episode")
	episodesCmd.Flags().Duration(common.ArgMinEpisodeDuration, 20*time.Minute, "Minimum runtime of an episode")
	episodesCmd.Flags().Duration(common.ArgShortChapterDuration, 30*time.Second, "Defines max length of short chapters")
	episodesCmd.Flags().String(common.ArgShortChapterMode, "none", "How to handle short chapters (none, discard, include)")
}
