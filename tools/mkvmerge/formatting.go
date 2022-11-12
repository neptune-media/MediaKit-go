package mkvmerge

import (
	"fmt"
	"path/filepath"
	"strings"
)

// FormatSplitOutputName is used to generate filenames when splitting
// a single input video into multiple output videos
func FormatSplitOutputName(outputFilename string, index int) string {
	dirname, filename := filepath.Split(outputFilename)
	ext := filepath.Ext(filename)
	basename := strings.TrimSuffix(filename, ext)
	name := fmt.Sprintf("%s-%03d%s", basename, index+1, ext)
	return filepath.Join(dirname, name)
}
