package cmd

import (
	"fmt"
	mediakit "github.com/neptune-media/MediaKit-go"
	"github.com/neptune-media/MediaKit-go/tasks"
	"github.com/neptune-media/MediaKit-go/tools"
	"github.com/neptune-media/MediaKit-go/tools/ffprobe"
	"github.com/neptune-media/MediaKit-go/tools/mkvmerge"
	"github.com/neptune-media/MediaKit-go/tools/mkvpropedit"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// splitCmd represents the split command
var splitCmd = &cobra.Command{
	Use:   "split [file]",
	Short: "Splits a multi-episode file into multiple files",
	Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := tools.NewChecks(
			ffprobe.NewCheck(),
			mkvmerge.NewCheck(),
			mkvpropedit.NewCheck(),
		).Run(); err != nil {
			return fmt.Errorf("pre-flight checks failed: %v", err)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		inputFilename := args[0]
		fmt.Printf("source file: %s\n", inputFilename)
		printCommands, err := cmd.Flags().GetBool("print")
		if err != nil {
			return fmt.Errorf("error while reading flag: %v", err)
		}

		iframesFilename, err := cmd.Flags().GetString("iframes")
		if err != nil {
			return fmt.Errorf("error while reading flag: %v", err)
		}

		opts := mediakit.EpisodeBuilderOptions{
			EndingChapterTime:    60 * time.Second,
			MinimumChapters:      2,
			MinimumEpisodeLength: 20 * time.Minute,
		}

		var frames []time.Duration
		if iframesFilename != "" {
			fmt.Printf("iframes file: %s\n", iframesFilename)
			iframesFile, err := os.Open(iframesFilename)
			if err != nil {
				return fmt.Errorf("error while reading IFrames: %v", err)
			}
			frames, err = mediakit.ReadFrameArray(iframesFile)
			if err != nil {
				return fmt.Errorf("error while reading IFrames: %v", err)
			}
		} else {
			frames, err = tasks.ReadVideoIFrames(inputFilename)
			if err != nil {
				return fmt.Errorf("error while reading IFrames: %v", err)
			}
		}
		opts.FrameSeeker = &mediakit.FrameSeeker{Frames: frames}

		// matroska-go outputs every block and is super noisy
		log.SetOutput(new(Sink))

		// Build episodes from file
		fmt.Printf("Building episodes from file\n")
		episodes, err := tasks.ReadVideoEpisodes(inputFilename, opts)
		if err != nil {
			return fmt.Errorf("error while reading episodes: %v", err)
		}
		fmt.Printf("Built %d episodes\n", len(episodes))

		runner := mkvmerge.NewSplitter(
			inputFilename,
			"output.mkv",
			episodes,
		)

		if printCommands {
			fmt.Fprintf(os.Stdout, "%s", runner.GetCommandString())
		} else {
			fmt.Printf("Splitting file into multiple episodes\n")
			if err := runner.Do(); err != nil {
				return fmt.Errorf("error while splitting file: %v\noutput from command:\n%s", err, runner.GetOutput())
			}

			fmt.Printf("Correcting episode chapter names\n")
			if err := mkvpropedit.FixEpisodeChapterNames(episodes, "output.mkv"); err != nil {
				return fmt.Errorf("error while writing chapters: %v", err)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(splitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// splitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// splitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	splitCmd.Flags().String("iframes", "", "Path to a file containing IFrame data for the video file")
	splitCmd.Flags().BoolP("print", "", false, "Print mkvmerge commands instead of running")
}
