package mediakit

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/go-logr/logr"

	"github.com/neptune-media/MediaKit-go/pkg/mediakit/utils/timing"
	"github.com/neptune-media/MediaKit-go/pkg/tools/ffprobe"
)

// FrameFilter returns true if a frame should be kept in the output
// when reading from parquet.
type FrameFilter func(frame ffprobe.Frame) bool

func FilterFrameByType(frameType string) FrameFilter {
	return func(frame ffprobe.Frame) bool {
		return frame.PictType == frameType
	}
}

var (
	FilterIFrames = FilterFrameByType("I")
)

func ReadFramesFromFilename(filename string, logger logr.Logger, filters ...FrameFilter) (ffprobe.FrameList, error) {
	closeAndLog := NewCloserLogger(logger)
	defer timing.NewFuncDurationPrinter(time.Now(), "finished reading frames")(logger)

	// Open parquet file for read
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error while opening frames file: %w", err)
	}
	defer closeAndLog(f, "error closing frames file")

	// Create new reader
	reader := NewParquetReader[ffprobe.Frame](
		f,
		WithLogger[ffprobe.Frame](logger),
		WithSchema[ffprobe.Frame](true),
	)
	defer closeAndLog(reader, "error closing parquet reader")

	// Read frames from parquet
	return ReadFramesFromParquet(reader, filters...)
}

func ReadFramesFromParquet(reader *ParquetReader[ffprobe.Frame], filters ...FrameFilter) (ffprobe.FrameList, error) {
	frames := make(ffprobe.FrameList, 0)
	rows := make(ffprobe.FrameList, 1000)

	for _, err := reader.Read(rows); err != io.EOF; _, err = reader.Read(rows) {
		for _, row := range rows {
			if !applyFiltersToFrame(row, filters...) {
				continue
			}
			frames = append(frames, row)
		}
	}

	return frames, nil
}

// applyFiltersToFrame checks the provided frame against a list of filters,
// returns true if the frame should be kept
func applyFiltersToFrame(frame ffprobe.Frame, filters ...FrameFilter) bool {
	for _, filter := range filters {
		if !filter(frame) {
			return false
		}
	}

	return true
}

func NewCloserLogger(logger logr.Logger) func(io.Closer, string) {
	return func(c io.Closer, errMsg string) {
		err := c.Close()
		if err != nil {
			logger.Error(err, errMsg)
		}
	}
}
