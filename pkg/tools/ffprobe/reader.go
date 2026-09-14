package ffprobe

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/go-viper/mapstructure/v2"

	"github.com/neptune-media/MediaKit-go/pkg/helpers"
)

type Reader struct {
	framesCh chan Frame

	statsLock sync.Mutex
	stats     ReaderStats
}

type ReaderStats struct {
	// Number of lines parsed into frames
	DecodedFrames uint64

	// Number of lines parsed into streams
	DecodedStreams uint64

	// Number of lines parsed
	ParsedLines uint64

	// Number of lines read
	ReadLines uint64

	// Stats generated time
	Time time.Time
}

func NewReader() *Reader {
	return &Reader{
		framesCh: make(chan Frame, 50),
	}
}

func (r *Reader) Close() error {
	if r.framesCh != nil {
		close(r.framesCh)
		r.framesCh = nil
	}

	return nil
}

func (r *Reader) Frames() <-chan Frame {
	return r.framesCh
}

func (r *Reader) Start(ctx context.Context, in io.Reader) error {
	if r.framesCh == nil {
		return fmt.Errorf("cannot start reader after closing")
	}

	var err error
	var frame Frame
	// var stream Stream
	var line string
	var parsed parsedLine
	reader := bufio.NewReader(in)

	r.statsLock.Lock()
	r.stats = ReaderStats{}
	r.statsLock.Unlock()

	for {
		line, err = readLine(reader)
		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}
		r.statsLock.Lock()
		r.stats.ReadLines++
		r.statsLock.Unlock()

		parsed = parseLine(line)
		r.statsLock.Lock()
		r.stats.ParsedLines++
		r.statsLock.Unlock()

		switch parsed.DataType {
		case "frame":
			frame, err = decodeParsedLine[Frame](parsed)
			if err != nil {
				close(r.framesCh)
				return err
			}

			r.statsLock.Lock()
			r.stats.DecodedFrames++
			r.statsLock.Unlock()

			r.framesCh <- frame
		case "stream":
			_, err = decodeParsedLine[Stream](parsed)
			if err != nil {
				close(r.framesCh)
				return err
			}

			r.statsLock.Lock()
			r.stats.DecodedStreams++
			r.statsLock.Unlock()
		}
	}

	return helpers.IgnoreEOF(err)
}

func (r *Reader) Stats() ReaderStats {
	r.statsLock.Lock()
	defer r.statsLock.Unlock()

	stats := r.stats
	stats.Time = time.Now()
	return stats
}

type parsedLine struct {
	DataType string
	Parts    map[string]string
}

func decodeParsedLine[T any](parsed parsedLine) (T, error) {
	var t T

	cfg := &mapstructure.DecoderConfig{
		Result:           &t,
		WeaklyTypedInput: true,
	}

	decoder, err := mapstructure.NewDecoder(cfg)
	if err != nil {
		return t, err
	}

	return t, decoder.Decode(parsed.Parts)
}

func parseLine(line string) parsedLine {
	raw := parsedLine{
		Parts: make(map[string]string),
	}

	parts := strings.Split(line, "|")
	raw.DataType, parts = parts[0], parts[1:]

	for _, part := range parts {
		p := strings.Split(part, "=")
		key, value := p[0], p[1]
		if value == "N/A" {
			continue
		}

		raw.Parts[key] = value
	}

	return raw
}

func readLine(r *bufio.Reader) (string, error) {
	// Read next line
	line, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}

	// Trim leading/trailing spaces
	line = strings.TrimSpace(line)
	return line, nil
}
