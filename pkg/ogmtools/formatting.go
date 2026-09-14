package ogmtools

import (
	"fmt"
	"time"
)

func ChapterNameString(id int) string {
	return fmt.Sprintf("CHAPTER%02dNAME=Chapter %02d", id+1, id+1)
}

func ChapterTimeString(id int, t time.Duration) string {
	return fmt.Sprintf("CHAPTER%02d=%s", id+1, FormatTimeString(t))
}

func FormatTimeString(t time.Duration) string {
	hours := t / time.Hour
	t %= time.Hour
	minutes := t / time.Minute
	t %= time.Minute
	seconds := float64(t) / float64(time.Second)

	// Generate a time code in format of HH:MM:SS.sss
	return fmt.Sprintf(
		"%02d:%02d:%06.3f",
		hours,
		minutes,
		seconds,
	)
}
