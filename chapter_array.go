package mediakit

import (
	"fmt"
	"io"
	"time"
)

type ChapterArray []Chapter

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
