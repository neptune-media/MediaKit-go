package mkvmerge

import (
	"fmt"
	"path/filepath"
	"strings"
)

func FormatSplitOutputName(outputFilename string, index int) string {
	dirname, filename := filepath.Split(outputFilename)
	ext := filepath.Ext(filename)
	basename := strings.TrimSuffix(filename, ext)
	name := fmt.Sprintf("%s-%03d%s", basename, index+1, ext)
	return filepath.Join(dirname, name)
}
