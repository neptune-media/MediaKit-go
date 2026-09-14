package mkvmerge

import (
	"github.com/neptune-media/MediaKit-go/pkg/tools"
)

func New() *tools.Executable {
	return &tools.Executable{Name: "mkvmerge"}
}
