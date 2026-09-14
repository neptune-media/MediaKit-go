package cmd

import "github.com/neptune-media/MediaKit-go/cmd/info"

func init() {
	rootCmd.AddCommand(info.Cmd)
}
