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
		saveFilename, err := cmd.Flags().GetString("save")
		if err != nil {
			fmt.Printf("error while reading flag: %v\n", err)
			return
		}
		fmt.Printf("Dumping iframes for %s\n", inputFilename)
		if saveFilename != "" {
			fmt.Printf("Saving iframes to %s\n", saveFilename)
		}

		if frames, err := mediakit.ReadVideoIFrames(inputFilename); err != nil {
			fmt.Printf("error while reading IFrames: %+v\n", err)
			return
		} else {
			if saveFilename != "" {
				f, err := os.Create(saveFilename)
				if err != nil {
					fmt.Printf("error while saving IFrames: %v\n", err)
					return
				}
				defer f.Close()

				if err := mediakit.FrameArray(frames).Write(f); err != nil {
					fmt.Printf("error while saving IFrames: %v\n", err)
					return
				}
			}
			for _, frame := range frames {
				fmt.Fprintf(os.Stdout, "%.03f s\n", frame.Seconds())
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
	iframesCmd.Flags().String("save", "", "Saves IFrames to specified file for later use")
}
