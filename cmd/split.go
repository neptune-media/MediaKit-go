package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
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
	"github.com/neptune-media/MediaKit-go/pkg/tools/mkvmerge"
	"github.com/neptune-media/MediaKit-go/pkg/tools/mkvpropedit"
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

		splitterBuilder := mkvmerge.NewCommandBuilder(mkvmerge.New(), zapr.NewLogger(logger))
		splitterBuilder = splitterBuilder.LowPriority().SplitEpisodes(episodes)
		tool := splitterBuilder.Build(cmd.Context(), inputFilename, "split.mkv")

		err = tool.Start()
		if err != nil {
			logger.Fatal("failed to start tool", zap.Error(err))
		}

		startTime := time.Now()
		logger.Info("splitting episodes")

		err = tool.Wait()
		if err != nil {
			logger.Fatal("failed to run tool", zap.Error(err))
		}
		stopTime := time.Now()

		duration := stopTime.Sub(startTime)
		logger.Info("finished splitting episodes", zap.Duration("duration", duration))

		startTime = time.Now()
		logger.Info("fixing chapter names")
		propeditBuilder := mkvpropedit.NewCommandBuilder(mkvpropedit.New(), zapr.NewLogger(logger))
		for idx, episode := range episodes {
			filename := mkvmerge.FormatSplitOutputName("split.mkv", idx)
			chFilename := fmt.Sprintf("%s.chapters", filename)
			tool := propeditBuilder.RenameChapters(chFilename).Build(cmd.Context(), filename)

			err = mediakit.WriteChapterNamesToFile(episode.Chapters, chFilename)
			if err != nil {
				logger.Fatal("failed to write chapter names", zap.Error(err))
			}

			err = tool.Start()
			if err != nil {
				logger.Fatal("failed to start tool", zap.Error(err))
			}
			err = tool.Wait()
			if err != nil {
				var toolErr *exec.ExitError
				if errors.As(err, &toolErr) {
					logger.Info("tool output", zap.String("stderr", string(toolErr.Stderr)))
				}
				logger.Fatal("failed to run tool", zap.Error(err))
			}
		}

		stopTime = time.Now()
		duration = stopTime.Sub(startTime)
		logger.Info("finished fixing chapter names", zap.Duration("duration", duration))

		// TODO: Finish this, but also, refactor this command and info/episodes, since
		// they are basically the exact same, one just keeps going and the other stops.
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
