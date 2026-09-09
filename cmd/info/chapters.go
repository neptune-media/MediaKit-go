package info

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-logr/zapr"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/neptune-media/MediaKit-go/cmd/common"
	"github.com/neptune-media/MediaKit-go/pkg/mediakit/utils"
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
		inputFilename := args[0]
		logger = logger.With(zap.String("job", filepath.Base(inputFilename)))
		logger.Info("using input file", zap.String("input-file", inputFilename))

		f, err := os.Open(inputFilename)
		if err != nil {
			logger.Fatal("failed to open file", zap.Error(err))
		}
		defer f.Close()

		chapters, err := utils.ReadVideoChapters(f, zapr.NewLogger(logger))
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
