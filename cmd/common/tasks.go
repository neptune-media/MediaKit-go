package common

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/go-logr/zapr"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/neptune-media/MediaKit-go/pkg/mediakit"
	"github.com/neptune-media/MediaKit-go/pkg/mediakit/builders/episode"
	"github.com/neptune-media/MediaKit-go/pkg/mediakit/utils"
	"github.com/neptune-media/MediaKit-go/pkg/tools/mkvmerge"
	"github.com/neptune-media/MediaKit-go/pkg/tools/mkvpropedit"
)

func defaultLogger(l *zap.Logger) *zap.Logger {
	if l == nil {
		return zap.NewNop()
	}

	return l
}

func BuildEpisodes(logger *zap.Logger, filename string) (mediakit.EpisodeList, error) {
	logger = defaultLogger(logger)

	// Get chapters
	chapters, err := GetVideoChapters(logger, filename)
	if err != nil {
		return nil, err
	}

	logger.Info("using input file", zap.String("input-file", filename))

	// Get frame data
	if viper.GetBool(ArgAlignChapters) {
		filename := viper.GetString(ArgFramesFile)
		logger.Info("using frames file", zap.String("frames-file", filename))
		frames, err := mediakit.ReadFramesFromFilename(filename, zapr.NewLogger(logger), mediakit.FilterIFrames)
		if err != nil {
			logger.Error("error while reading frame data", zap.Error(err))
			return nil, err
		}

		logger.Info("aligning chapters to I-frames")
		seeker := episode.NewFrameSeeker(frames)
		aligner := episode.Aligner{
			Logger: zapr.NewLogger(logger),
		}

		chapters, err = aligner.AlignChaptersToIFrames(chapters, seeker)
		if err != nil {
			logger.Error("error while aligning chapters to I-frames", zap.Error(err))
			return nil, err
		}
	}
	totalChapterRuntime := chapters.Runtime()

	logger.Info("building episodes from chapters")
	builder := episode.NewEpisodeBuilder(
		episode.WithLogger(zapr.NewLogger(logger)),
		episode.WithOptions(NewEpisodeOptionsBuilder().Build()),
	)
	episodes, err := builder.WithChapters(chapters).BuildAll()
	if err != nil {
		logger.Error("error while building episodes", zap.Error(err))
		return nil, err
	}
	totalEpisodeRuntime := episodes.Runtime()
	logger.Info("episodes runtime",
		zap.Duration("pre-split", totalChapterRuntime),
		zap.Duration("post-split", totalEpisodeRuntime))

	return episodes, nil
}

func FixChapterNames(ctx context.Context, logger *zap.Logger, splitFilename string, episodes mediakit.EpisodeList) error {
	logger = defaultLogger(logger)
	var err error

	startTime := time.Now()
	logger.Info("fixing chapter names")
	builder := mkvpropedit.NewCommandBuilder(mkvpropedit.New(), zapr.NewLogger(logger))

	for idx, episode := range episodes {
		filename := mkvmerge.FormatSplitOutputName("split.mkv", idx)
		chFilename := fmt.Sprintf("%s.chapters", filename)
		tool := builder.RenameChapters(chFilename).Build(ctx, filename)

		err = mediakit.WriteChapterNamesToFile(episode.Chapters, chFilename)
		if err != nil {
			logger.Error("error while writing chapter names", zap.Error(err))
			return err
		}

		err = tool.Start()
		if err != nil {
			logger.Error("error while starting tool", zap.Error(err))
		}
		err = tool.Wait()
		if err != nil {
			var toolErr *exec.ExitError
			if errors.As(err, &toolErr) {
				logger.Info("tool output", zap.String("stderr", string(toolErr.Stderr)))
			}

			logger.Error("error running tool", zap.Error(err))
			return err
		}
	}

	stopTime := time.Now()
	duration := stopTime.Sub(startTime)
	logger.Info("finished fixing chapter names", zap.Duration("duration", duration))
	return nil
}

func GetVideoChapters(logger *zap.Logger, filename string) (mediakit.ChapterList, error) {
	logger = defaultLogger(logger)
	logger.Info("using input file", zap.String("input-file", filename))

	// Open input for read
	f, err := os.Open(filename)
	if err != nil {
		logger.Error("error while opening file", zap.Error(err))
		return nil, err
	}
	defer f.Close()

	return utils.ReadVideoChapters(f, zapr.NewLogger(logger))
}

func SplitEpisodes(ctx context.Context, logger *zap.Logger, filename string, episodes mediakit.EpisodeList) error {
	logger = defaultLogger(logger)

	splitter := mkvmerge.NewCommandBuilder(mkvmerge.New(), zapr.NewLogger(logger))
	if viper.GetBool(ArgLowPriority) {
		splitter = splitter.LowPriority()
	}
	splitter = splitter.SplitEpisodes(episodes)

	tool := splitter.Build(ctx, filename, "split.mkv")
	err := tool.Start()
	if err != nil {
		logger.Error("error starting tool", zap.Error(err))
		return err
	}

	startTime := time.Now()
	logger.Info("splitting episodes")

	err = tool.Wait()
	if err != nil {
		logger.Error("error running tool", zap.Error(err))
		return err
	}
	stopTime := time.Now()

	duration := stopTime.Sub(startTime)
	logger.Info("finished splitting episodes", zap.Duration("duration", duration))
	return nil
}
