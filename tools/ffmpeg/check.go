package ffmpeg

import (
	"github.com/neptune-media/MediaKit-go/tools"
)

// NewCheck returns a Check configured for checking if ffmpeg can be run
func NewCheck() tools.Check {
	return &tools.ExecutableCheck{Executable: "ffmpeg"}
}
