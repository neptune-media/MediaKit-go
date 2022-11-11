package cmd

import (
	"fmt"
	mediakit "github.com/neptune-media/MediaKit-go"
	"github.com/neptune-media/MediaKit-go/tasks"
	"github.com/neptune-media/MediaKit-go/tools"
	"github.com/neptune-media/MediaKit-go/tools/ffprobe"
	"github.com/neptune-media/MediaKit-go/tools/mkvmerge"
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
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := tools.NewChecks(ffprobe.NewCheck(), mkvmerge.NewCheck()).Run(); err != nil {
			return fmt.Errorf("pre-flight checks failed: %v", err)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		inputFilename := args[0]
		alignChapters, err := cmd.Flags().GetBool("align-chapters")
		if err != nil {
			return fmt.Errorf("error while reading flag: %v", err)
		}
		iframesFilename, err := cmd.Flags().GetString("iframes")
		if err != nil {
			return fmt.Errorf("error while reading flag: %v", err)
		}

		opts := newEpisodeBuilderOptionsFromFlags(cmd)

		if alignChapters {
			var frames []time.Duration
			frames, err = loadIFrames(inputFilename, iframesFilename)
			if err != nil {
				return fmt.Errorf("error while reading IFrames: %v", err)
			}
			opts.FrameSeeker = &mediakit.FrameSeeker{Frames: frames}
		}

		fmt.Printf("Dumping episodes for %s\n", inputFilename)

		// matroska-go outputs every block and is super noisy
		log.SetOutput(new(Sink))

		if episodes, err := tasks.ReadVideoEpisodes(inputFilename, opts); err != nil {
			return fmt.Errorf("error while reading episodes: %v", err)
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

		return nil
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
	addEpisodeBuilderFlags(episodesCmd)
}

func addEpisodeBuilderFlags(cmd *cobra.Command) {
	cmd.Flags().Duration("ending-chapter-time", 60, "Chapters longer than this will continue the episode")
	cmd.Flags().Int("minimum-chapters", 2, "Minimum number of chapters in an episode")
	cmd.Flags().Duration("minimum-episode-length", 20, "Minimum runtime of an episode")
}

func newEpisodeBuilderOptionsFromFlags(cmd *cobra.Command) mediakit.EpisodeBuilderOptions {
	endingChapterTime, _ := cmd.Flags().GetDuration("ending-chapter-time")
	minimumChapters, _ := cmd.Flags().GetInt("minimum-chapters")
	minimumEpisodeLength, _ := cmd.Flags().GetDuration("minimum-episode-length")
	return mediakit.EpisodeBuilderOptions{
		EndingChapterTime:    endingChapterTime * time.Second,
		MinimumChapters:      minimumChapters,
		MinimumEpisodeLength: minimumEpisodeLength * time.Minute,
	}
}

func loadIFrames(sourceFilename, iframesFilename string) ([]time.Duration, error) {
	if iframesFilename == "" {
		return tasks.ReadVideoIFramesFromFile(sourceFilename)
	}

	fmt.Printf("iframes file: %s\n", iframesFilename)
	f, err := os.Open(iframesFilename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return mediakit.ReadFrameArray(f)
}
