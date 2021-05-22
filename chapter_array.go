package mediakit

import (
	"fmt"
	"github.com/neptune-media/MediaKit-go/ogmtools"
	"io"
)

type ChapterArray []Chapter

func (a ChapterArray) WriteTo(w io.Writer) (int64, error) {
	var n int64 = 0
	offset := []Chapter(a)[0].StartTime()
	for i, chapter := range []Chapter(a) {
		// Offset chapter start time by first chapter to get relative start time
		startTime := chapter.StartTime() - offset

		// Write chapter timecode
		count, err := fmt.Fprintln(w, ogmtools.ChapterTimeString(i, startTime))
		n += int64(count)
		if err != nil {
			return n, err
		}

		// Write chapter name
		count, err = fmt.Fprintln(w, ogmtools.ChapterNameString(i))
		n += int64(count)
		if err != nil {
			return n, err
		}
	}
	return n, nil
}
