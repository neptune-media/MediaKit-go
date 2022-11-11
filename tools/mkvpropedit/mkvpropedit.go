package mkvpropedit

import (
	"bytes"
	"context"
	"fmt"
	"github.com/neptune-media/MediaKit-go"
	"github.com/neptune-media/MediaKit-go/tools"
	"github.com/neptune-media/MediaKit-go/tools/mkvmerge"
	"io"
	"os"
)

type MKVPropedit struct {
	ChaptersFilename string
	Filename         string
	LowPriority      bool

	stdout       []byte
	stderr       []byte
	stdoutBuffer bytes.Buffer
	stderrBuffer bytes.Buffer
}

func (m *MKVPropedit) Do() error {
	return m.DoWithContext(context.Background())
}

func (m *MKVPropedit) DoWithContext(ctx context.Context) error {
	// Reset buffers
	m.stdout = make([]byte, 0)
	m.stderr = make([]byte, 0)
	m.stdoutBuffer.Reset()
	m.stderrBuffer.Reset()

	// Execute
	err := tools.ExecTool(ctx, m)

	// Copy output to buffer for later
	m.stdout = make([]byte, m.stdoutBuffer.Len())
	m.stderr = make([]byte, m.stderrBuffer.Len())
	copy(m.stdout, m.stdoutBuffer.Bytes())
	copy(m.stderr, m.stderrBuffer.Bytes())

	return err
}

func (m *MKVPropedit) GetCommand() string {
	return "mkvpropedit"
}

func (m *MKVPropedit) GetCommandArgs() []string {
	return []string{
		m.Filename,
		"-c",
		m.ChaptersFilename,
	}
}

func (m *MKVPropedit) GetStdout() []byte {
	return m.stdout
}

func (m *MKVPropedit) GetStderr() []byte {
	return m.stderr
}

func (m *MKVPropedit) GetOutputBuffers() (io.Writer, io.Writer) {
	return &m.stdoutBuffer, &m.stderrBuffer
}

func (m *MKVPropedit) IsLowPriority() bool {
	return m.LowPriority
}

func FixEpisodeChapterNames(episodes []mediakit.Episode, filename string) error {
	for i, episode := range episodes {
		videoFilename := mkvmerge.FormatSplitOutputName(filename, i)
		chFilename := fmt.Sprintf("%s.chapters", videoFilename)
		err := writeChapterNamesToFile(episode.Chapters, chFilename)
		if err != nil {
			return err
		}

		runner := &MKVPropedit{
			ChaptersFilename: chFilename,
			Filename:         videoFilename,
		}
		if err := runner.Do(); err != nil {
			fmt.Printf("error while updating file: %v\n", err)
			fmt.Printf("output from command:\n%s\n%s\n", runner.GetStdout(), runner.GetStderr())
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
