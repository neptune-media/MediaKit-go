package cmd

import (
	"fmt"
	"github.com/neptune-media/MediaKit-go"
	"github.com/spf13/cobra"
	"log"
	"os"
)

// Sink is a dummy writer that doesn't do anything
type Sink int

func (s Sink) Write(p []byte) (int, error) {
	return len(p), nil
}

// chaptersCmd represents the chapters command
var chaptersCmd = &cobra.Command{
	Use:   "chapters [file]",
	Short: "Prints a list of chapters in a given file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		inputFilename := args[0]
		fmt.Printf("Dumping chapters for %s\n", inputFilename)

		// matroska-go outputs every block and is super noisy
		log.SetOutput(new(Sink))
		if chapters, err := mediakit.ReadVideoChapters(inputFilename); err != nil {
			fmt.Printf("error while reading chapters: %+v\n", err)
			return
		} else {
			for i, chapter := range chapters {
				if chapter.Enabled {
					fmt.Fprintf(os.Stdout,
						"%d %d - %d (%.1f seconds)\n",
						i,
						chapter.TimeStart,
						chapter.TimeEnd,
						chapter.Runtime().Seconds(),
					)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(chaptersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// chaptersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// chaptersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
