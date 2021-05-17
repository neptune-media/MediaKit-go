package cmd

import (
	"fmt"
	mediakit "github.com/neptune-media/MediaKit-go"
	"os"

	"github.com/spf13/cobra"
)

// iframesCmd represents the iframes command
var iframesCmd = &cobra.Command{
	Use:   "iframes [file]",
	Short: "Prints a list of frame numbers of I-frames in the given file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		inputFilename := args[0]
		fmt.Printf("Dumping iframes for %s\n", inputFilename)

		if frames, err := mediakit.ReadVideoIFrames(inputFilename); err != nil {
			fmt.Printf("error while reading IFrames: %+v\n", err)
			return
		} else {
			for _, frame := range frames {
				fmt.Fprintf(os.Stdout, "%d\n", frame)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(iframesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// iframesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// iframesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
