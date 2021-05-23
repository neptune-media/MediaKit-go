package tasks

import (
	"github.com/neptune-media/MediaKit-go"
	"github.com/pixelbender/go-matroska/matroska"
)

func ReadVideoChapters(filename string) ([]mediakit.Chapter, error) {
	chapters := make([]mediakit.Chapter, 0)

	f, err := matroska.Decode(filename)
	if err != nil {
		return nil, err
	}

	timeScale := f.Segment.Info[0].TimecodeScale

	for _, chapter := range f.Segment.Chapters {
		for _, atom := range chapter.Atoms {
			chapters = append(chapters, mediakit.Chapter{
				ID:        uint64(atom.ID),
				TimeStart: int64(atom.TimeStart) / int64(timeScale),
				TimeEnd:   int64(atom.TimeEnd) / int64(timeScale),
				Enabled:   atom.Enabled,
			})
		}
	}

	return chapters, nil
}
