package episode

import (
	"time"

	"github.com/neptune-media/MediaKit-go/pkg/tools/ffprobe"
)

type FrameSeeker interface {
	// Current returns the time position of the current frame
	Current() time.Duration

	// EOF returns true when the seeker has reached the end of frames
	EOF() bool

	// Next advances the seeker to the next frame and returns the new frame time position
	Next() time.Duration

	// Peek returns the time position of the next frame without advancing the seeker
	Peek() time.Duration

	// Reset moves the seeker back to the start of the frame list
	Reset()
}

type frameSeeker struct {
	Frames ffprobe.FrameList
	pos    int
}

func NewFrameSeeker(frames ffprobe.FrameList) FrameSeeker {
	return &frameSeeker{
		Frames: frames,
	}
}

func (f *frameSeeker) Current() time.Duration {
	return f.Frames[f.pos].Timecode()
}

func (f *frameSeeker) EOF() bool {
	return !(f.pos < len(f.Frames))
}

func (f *frameSeeker) Next() time.Duration {
	if !f.EOF() {
		f.pos += 1
	}

	return f.Frames[f.pos].Timecode()
}

func (f *frameSeeker) Peek() time.Duration {
	if (f.pos + 1) < len(f.Frames) {
		return f.Frames[f.pos+1].Timecode()
	}

	return f.Frames[f.pos].Timecode()
}

func (f *frameSeeker) Reset() {
	f.pos = 0
}
