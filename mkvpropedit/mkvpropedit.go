package mkvpropedit

import (
	"fmt"
	"github.com/neptune-media/MediaKit-go"
	"github.com/neptune-media/MediaKit-go/mkvmerge"
	"os"
	"os/exec"
)

type Runner struct {
	ChaptersFilename string
	Filename         string

	output []byte
}

func (r *Runner) Do() error {
	c := exec.Command(
		"mkvpropedit",
		r.Filename,
		"-c",
		r.ChaptersFilename,
	)
	o, err := c.CombinedOutput()

	r.output = make([]byte, len(o))
	copy(r.output, o)
	return err
}

func (r *Runner) Output() []byte {
	o := make([]byte, len(r.output))
	copy(o, r.output)
	return o
}

func FixEpisodeChapterNames(episodes []mediakit.Episode, filename string) error {
	for i, episode := range episodes {
		videoFilename := mkvmerge.FormatSplitOutputName(filename, i)
		chFilename := fmt.Sprintf("%s.chapters", videoFilename)
		err := writeChapterNamesToFile(episode.Chapters, chFilename)
		if err != nil {
			return err
		}

		runner := &Runner{
			ChaptersFilename: chFilename,
			Filename:         videoFilename,
		}
		if err := runner.Do(); err != nil {
			fmt.Printf("error while updating file: %v\n", err)
			fmt.Printf("output from command:\n%s\n", runner.Output())
			return err
		}

		if err := os.Remove(chFilename); err != nil {
			return err
		}
	}

	return nil
}

func writeChapterNamesToFile(chapters mediakit.ChapterArray, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := chapters.WriteTo(f); err != nil {
		return err
	}

	return nil
}
