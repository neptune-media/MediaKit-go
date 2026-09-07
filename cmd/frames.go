package cmd

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/parquet-go/parquet-go/compress/snappy"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/neptune-media/MediaKit-go/pkg/mediakit"
	"github.com/neptune-media/MediaKit-go/pkg/tools/ffprobe"
)

// framesCmd represents the frames command
var framesCmd = &cobra.Command{
	Use:   "frames [file]",
	Short: "Extracts a list of frames from the given file to parquet",
	Long: `Extracts a list of frames from the given file and writes the
information to parquet.  Can optionally print a list of frames to stdout
instead.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		tool := ffprobe.New()
		return tool.Validate(cmd.Context())
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Setup logging
		logger := initLogger()
		defer logger.Sync()

		inputFilename := args[0]
		logger = logger.With(zap.String("job", filepath.Base(inputFilename)))
		logger.Info("using input file", zap.String("input-file", inputFilename))

		framesFilename := viper.GetString(ArgFramesFile)
		logger.Info("saving frame info to file", zap.String("frames-file", framesFilename))

		outFile, err := os.Create(framesFilename)
		if err != nil {
			logger.Fatal("failed to create output file", zap.Error(err))
		}
		defer outFile.Close()

		writer := mediakit.NewParquetWriter[ffprobe.Frame](
			outFile,
			mediakit.WithCompression[ffprobe.Frame](new(snappy.Codec)),
		)
		defer writer.Close()

		builder := ffprobe.NewBuilder(ffprobe.New(), logger, inputFilename)
		builder = builder.GetFramesCount().GetFrames()
		if viper.GetBool(ArgThreads) {
			builder = builder.UseThreads(0)
		}

		if viper.GetBool(ArgLowPriority) {
			builder = builder.UseLowPriority()
		}

		tool := builder.Build(cmd.Context())
		stdout, err := tool.StdoutPipe()
		if err != nil {
			logger.Fatal("failed to get stdout pipe", zap.Error(err))
		}
		defer stdout.Close()

		// stderr, err := tool.StderrPipe()
		// if err != nil {
		// 	logger.Fatal("failed to get stderr pipe", zap.Error(err))
		// }
		// defer stderr.Close()

		err = tool.Start()
		if err != nil {
			logger.Fatal("failed to start tool", zap.Error(err))
		}

		startTime := time.Now()
		logger.Info("dumping frames")

		frameReader := ffprobe.NewReader()
		frames := frameReader.Frames()
		go func(r *ffprobe.Reader) {
			err := r.Start(cmd.Context(), stdout)
			if err != nil {
				logger.Fatal("failed to read frame", zap.Error(err))
			}
			defer r.Close()
		}(frameReader)

		// Stats printer
		statsCancelFn := newStatsPrinter(cmd.Context(), logger, viper.GetDuration(ArgStatsInterval), frameReader)
		// statsCtx, statsCancelFn := context.WithCancel(cmd.Context())
		// go func(ctx context.Context, interval time.Duration) {
		// 	ticker := time.NewTicker(interval)
		// 	lastStats := ffprobe.ReaderStats{
		// 		Time: time.Now(),
		// 	}
		// 	for {
		// 		select {
		// 		case <-ctx.Done():
		// 			return
		// 		case t := <-ticker.C:
		// 			stats := frameReader.Stats()
		// 			deltaFrames := stats.DecodedFrames - lastStats.DecodedFrames
		// 			rate := float64(deltaFrames) / t.Sub(lastStats.Time).Seconds()
		// 			logger.Info("read progress", zap.Uint("decoded-frames", uint(stats.DecodedFrames)), zap.Uint("fps", uint(rate)))
		// 		}
		// 	}
		// }(statsCtx, time.Minute)

		err = handleItems(writer, frames)
		if err != nil {
			logger.Fatal("failed to write frames", zap.Error(err))
		}

		err = tool.Wait()
		if err != nil {
			// stdoutO, _ := io.ReadAll(stdout)
			// stderrO, _ := io.ReadAll(stderr)
			// if err != nil {
			// 	logger.Fatal("error while reading stderr", zap.Error(err))
			// }
			// logger.Error("error while running tool", zap.String("stderr", string(stderrO)), zap.String("stdout", string(stdoutO)))
			logger.Fatal("failed to run tool", zap.Error(err))
		}
		stopTime := time.Now()

		statsCancelFn()
		duration := stopTime.Sub(startTime)
		logger.Info("finished dumping iframes", zap.Duration("duration", duration))
	},
}

func init() {
	rootCmd.AddCommand(framesCmd)

	framesCmd.Flags().String(ArgFramesFile, "", "path to save frame information to")
	framesCmd.Flags().Bool(ArgLowPriority, false, "When set, runs subprocesses at a lower priority")
	framesCmd.Flags().Duration(ArgStatsInterval, time.Minute, "Specifies the interval for printing out stats")
	framesCmd.Flags().Bool(ArgThreads, false, "When set, set subprocess thread flags when appropriate")
	framesCmd.MarkFlagRequired(ArgFramesFile)
}

func handleItems[T any](writer *mediakit.ParquetWriter[T], items <-chan T) error {
	var err error

	for {
		select {
		case item, ok := <-items:
			if !ok {
				return nil
			}

			_, err = writer.Write([]T{item})
			if err != nil {
				return err
			}
		}
	}
}

func newStatsPrinter(ctx context.Context, logger *zap.Logger, interval time.Duration, reader *ffprobe.Reader) context.CancelFunc {
	ctx, cancelFn := context.WithCancel(ctx)

	f := func() {
		ticker := time.NewTicker(interval)
		lastStats := ffprobe.ReaderStats{
			Time: time.Now(),
		}

		for {
			select {
			case <-ctx.Done():
				return
			case t := <-ticker.C:
				stats := reader.Stats()
				deltaFrames := stats.DecodedFrames - lastStats.DecodedFrames
				rate := float64(deltaFrames) / t.Sub(lastStats.Time).Seconds()
				lastStats = stats
				logger.Info("read progress", zap.Uint("decoded-frames", uint(stats.DecodedFrames)), zap.Uint("decoded-streams", uint(stats.DecodedStreams)), zap.Uint("fps", uint(rate)))
			}
		}
	}

	go f()
	return cancelFn
}
