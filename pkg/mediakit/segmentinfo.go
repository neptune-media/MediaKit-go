package mediakit

import "time"

type SegmentInfo struct {
	// Timescale is the multiplier to apply to chapters
	// to get the actual time.
	Timescale time.Duration
}
