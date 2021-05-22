package mkvmerge

import (
	"fmt"
	mediakit "github.com/neptune-media/MediaKit-go"
	"strings"
)

func NewSplitter(input, output string, episodes []mediakit.Episode) *Runner {

	// Generate split part list
	splitPartList := make([]string, len(episodes))
	for i, episode := range episodes {
		splitPartList[i] = fmt.Sprintf(
			"%dms-%dms",
			episode.Chapters[0].TimeStart,
			episode.Chapters[len(episode.Chapters)-1].TimeEnd,
		)
	}

	args := []string{
		"--split",
		fmt.Sprintf("parts:%s", strings.Join(splitPartList, ",")),
	}

	return &Runner{
		Args:           args,
		InputFilename:  input,
		OutputFilename: output,
	}
}
