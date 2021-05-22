package mediakit

import "time"

type FrameSeeker struct {
	Frames []time.Duration
	pos    int
}

func (f *FrameSeeker) AtEnd() bool {
	return !(f.pos < len(f.Frames))
}

func (f *FrameSeeker) Current() time.Duration {
	return f.Frames[f.pos]
}

func (f *FrameSeeker) Next() {
	if !f.AtEnd() {
		f.pos += 1
	}
}

func (f *FrameSeeker) Peek() time.Duration {
	if f.pos+1 < len(f.Frames) {
		return f.Frames[f.pos+1]
	}

	return f.Frames[f.pos]
}

func (f *FrameSeeker) Reset() {
	f.pos = 0
}
