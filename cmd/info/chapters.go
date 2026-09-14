package info

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/neptune-media/MediaKit-go/cmd/common"
)

// chaptersCmd represents the chapters command
var chaptersCmd = &cobra.Command{
	Use:   "chapters [file]",
	Short: "Prints a list of chapters in a given file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Setup logging
		logger := common.InitLogger()
		defer logger.Sync()

		// Get input
		filename := args[0]
		logger = logger.With(zap.String("job", filepath.Base(filename)))
		chapters, err := common.GetVideoChapters(logger, filename)
		if err != nil {
			logger.Fatal("failed to parse file", zap.Error(err))
		}

		fmt.Println("Chapters")
		fmt.Println()
		for _, ch := range chapters {
			fmt.Printf("Ch %2d\t%8.1f\t%8.1f (%s)\t%s\n", ch.ID, ch.StartTime().Seconds(), ch.EndTime().Seconds(), ch.Runtime().String(), ch.Title)
		}
	},
}

func init() {
	infoCmd.AddCommand(chaptersCmd)
}
