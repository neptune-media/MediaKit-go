package cmd

import (
	"fmt"
	mediakit "github.com/neptune-media/MediaKit-go"
	"github.com/neptune-media/MediaKit-go/tasks"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// episodesCmd represents the episodes command
var episodesCmd = &cobra.Command{
	Use:   "episodes [file]",
	Short: "Scans chapters in a file and combines them into episodes",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		inputFilename := args[0]
		alignChapters, err := cmd.Flags().GetBool("align-chapters")
		if err != nil {
			fmt.Printf("error while reading flag: %v\n", err)
			return
		}
		iframesFilename, err := cmd.Flags().GetString("iframes")
		if err != nil {
			fmt.Printf("error while reading flag: %v\n", err)
			return
		}

		opts := mediakit.EpisodeBuilderOptions{
			EndingChapterTime:    60 * time.Second,
			MinimumChapters:      2,
			MinimumEpisodeLength: 20 * time.Minute,
		}

		if alignChapters {
			var frames []time.Duration
			if iframesFilename != "" {
				fmt.Printf("iframes file: %s\n", iframesFilename)
				iframesFile, err := os.Open(iframesFilename)
				if err != nil {
					fmt.Printf("error while reading IFrames: %+v\n", err)
					return
				}
				frames, err = mediakit.ReadFrameArray(iframesFile)
				if err != nil {
					fmt.Printf("error while reading IFrames: %+v\n", err)
					return
				}
			} else {
				frames, err = tasks.ReadVideoIFrames(inputFilename)
				if err != nil {
					fmt.Printf("error while reading IFrames: %+v\n", err)
					return
				}
			}
			opts.FrameSeeker = &mediakit.FrameSeeker{Frames: frames}
		}

		fmt.Printf("Dumping episodes for %s\n", inputFilename)

		// matroska-go outputs every block and is super noisy
		log.SetOutput(new(Sink))

		if episodes, err := tasks.ReadVideoEpisodes(inputFilename, opts); err != nil {
			fmt.Printf("error while reading episodes: %+v\n", err)
			return
		} else {
			for i, episode := range episodes {
				fmt.Fprintf(os.Stdout,
					"%d %d - %d (%.1f minutes)\n",
					i,
					episode.Chapters[0].TimeStart,
					episode.Chapters[len(episode.Chapters)-1].TimeEnd,
					episode.Runtime().Minutes(),
				)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(episodesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// episodesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// episodesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	episodesCmd.Flags().BoolP("align-chapters", "", false, "Align chapter markers to iframes")
	episodesCmd.Flags().String("iframes", "", "Path to a file containing IFrame data for the video file")
}
