package mediakit

import (
	"fmt"
	"github.com/pixelbender/go-matroska/matroska"
	"io"
	"time"
)

type Chapter struct {
	ID        uint64
	TimeStart int64
	TimeEnd   int64
	Enabled   bool
}

type ChapterArray []Chapter

func ReadVideoChapters(filename string) ([]Chapter, error) {
	chapters := make([]Chapter, 0)

	f, err := matroska.Decode(filename)
	if err != nil {
		return nil, err
	}

	timeScale := f.Segment.Info[0].TimecodeScale

	for _, chapter := range f.Segment.Chapters {
		for _, atom := range chapter.Atoms {
			chapters = append(chapters, Chapter{
				ID:        uint64(atom.ID),
				TimeStart: int64(atom.TimeStart) / int64(timeScale),
				TimeEnd:   int64(atom.TimeEnd) / int64(timeScale),
				Enabled:   atom.Enabled,
			})
		}
	}

	return chapters, nil
}

func (c Chapter) Runtime() time.Duration {
	return c.EndTime() - c.StartTime()
}

func (c Chapter) EndTime() time.Duration {
	return time.Duration(c.TimeEnd) * time.Millisecond
}

func (c Chapter) StartTime() time.Duration {
	return time.Duration(c.TimeStart) * time.Millisecond
}

func (a ChapterArray) WriteTo(w io.Writer) (int64, error) {
	var n int64 = 0
	offset := []Chapter(a)[0].StartTime()
	for i, chapter := range []Chapter(a) {
		// Offset chapter start time by first chapter to get relative start time
		chStartTime := chapter.StartTime() - offset
		hours := chStartTime / time.Hour
		chStartTime %= time.Hour
		minutes := chStartTime / time.Minute
		chStartTime %= time.Minute
		seconds := float64(chStartTime) / float64(time.Second)

		// Generate a time code in format of HH:MM:SS.sss
		chTimeStr := fmt.Sprintf(
			"%02d:%02d:%06.3f",
			hours,
			minutes,
			seconds,
		)

		// Write chapter timecode
		count, err := fmt.Fprintf(w, "CHAPTER%02d=%s\n", i+1, chTimeStr)
		n += int64(count)
		if err != nil {
			return n, err
		}

		// Write chapter name
		count, err = fmt.Fprintf(w, "CHAPTER%02dNAME=Chapter %02d\n", i+1, i+1)
		n += int64(count)
		if err != nil {
			return n, err
		}
	}
	return n, nil
}
