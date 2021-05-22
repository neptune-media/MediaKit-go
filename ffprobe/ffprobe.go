package ffprobe

import (
	"bytes"
	"io"
	"os/exec"
	"sync"
)

type Runner struct {
	Filename string
}

func (r *Runner) GetFrames() (*Output, error) {
	data, err := r.execAndWait("-select_streams",
		"v",
		"-show_frames",
		"-of",
		"json=compact=1",
		r.Filename)
	if err != nil {
		return nil, err
	}

	return ReadFFProbeOutput(bytes.NewReader(data))
}

func (r *Runner) execAndWait(args ...string) ([]byte, error) {
	c := exec.Command("ffprobe", args...)

	stdout, err := c.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := c.Start(); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	bufLock := &sync.Mutex{}
	go func() {
		bufLock.Lock()
		defer bufLock.Unlock()
		for {
			if _, err := io.Copy(&buf, stdout); err != nil {
				return
			}
		}
	}()

	if err := c.Wait(); err != nil {
		return nil, err
	}
	bufLock.Lock()
	defer bufLock.Unlock()

	arr := make([]byte, buf.Len())
	copy(arr, buf.Bytes())

	return arr, nil
}
