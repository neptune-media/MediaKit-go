package mediakit

import (
	"github.com/pixelbender/go-matroska/matroska"
	"time"
)

type Chapter struct {
	ID        uint64
	TimeStart int64
	TimeEnd   int64
	Enabled   bool
}

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
	return time.Duration(c.TimeEnd-c.TimeStart) * time.Millisecond
}
