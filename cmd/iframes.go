package cmd

import (
	"fmt"
	mediakit "github.com/neptune-media/MediaKit-go"
	"github.com/neptune-media/MediaKit-go/tasks"
	"github.com/neptune-media/MediaKit-go/tools"
	"github.com/neptune-media/MediaKit-go/tools/ffprobe"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// iframesCmd represents the iframes command
var iframesCmd = &cobra.Command{
	Use:   "iframes [file]",
	Short: "Prints a list of frame numbers of I-frames in the given file",
	Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := tools.NewChecks(ffprobe.NewCheck()).Run(); err != nil {
			return fmt.Errorf("pre-flight checks failed: %v", err)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		inputFilename := args[0]
		saveFilename, err := cmd.Flags().GetString("save")
		if err != nil {
			return fmt.Errorf("error while reading flag: %v", err)

		}
		fmt.Printf("Dumping iframes for %s\n", inputFilename)
		if saveFilename != "" {
			fmt.Printf("Saving iframes to %s\n", saveFilename)
		}

		startTime := time.Now()
		if frames, err := tasks.ReadVideoIFrames(inputFilename); err != nil {
			return fmt.Errorf("error while reading IFrames: %v", err)
		} else {
			readTime := time.Now().Sub(startTime)
			fmt.Printf("Took %s to analyze file\n", readTime.String())
			if saveFilename != "" {
				f, err := os.Create(saveFilename)
				if err != nil {
					return fmt.Errorf("error while saving IFrames: %v", err)
				}
				defer f.Close()

				if _, err := mediakit.FrameArray(frames).WriteTo(f); err != nil {
					return fmt.Errorf("error while saving IFrames: %v", err)
				}
			} else {
				for _, frame := range frames {
					fmt.Fprintf(os.Stdout, "%.03f s\n", frame.Seconds())
				}
			}
		}

		return nil
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
