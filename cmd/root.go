package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mediakit",
	Short: "Sample apps for testing library functionality",
	Long: `MediaKit provides a library for extracting video chunks from
a single video track, using ffmpeg.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return bindFlagsToViper(cmd)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func bindFlagsToViper(cmd *cobra.Command) error {
	flagSets := []*pflag.FlagSet{
		cmd.PersistentFlags(),
		cmd.Flags(),
	}
	for _, flags := range flagSets {
		if err := viper.BindPFlags(flags); err != nil {
			return fmt.Errorf("error while binding flags: %s", err)
		}
	}

	return nil
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().String(ArgLogFormat, "console", "Log format (json, console)")
	rootCmd.PersistentFlags().String(ArgLogLevel, "info", "Log level (debug, info, warn, error, fatal)")
	rootCmd.PersistentFlags().Bool(ArgLowPriority, false, "Runs subprocesses (codec/mkvmerge/etc) at a lower process priority")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetEnvPrefix("mediakit")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv() // read in environment variables that match
}

func initLogger() *zap.Logger {
	loggerCfg := newLoggerConfig()
	logger, err := loggerCfg.Build()
	if err != nil {
		panic(err)
	}
	return logger
}

func newLoggerConfig() *zap.Config {
	level, err := zap.ParseAtomicLevel(viper.GetString(ArgLogLevel))
	if err != nil {
		panic(err)
	}

	cfg := &zap.Config{
		Development:      false,
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		ErrorOutputPaths: []string{"stderr"},
		Level:            level,
		OutputPaths:      []string{"stdout"},
	}

	switch viper.GetString(ArgLogFormat) {
	case "json":
		cfg.Encoding = "json"
	default:
		cfg.Encoding = "console"
		cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		cfg.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
		cfg.EncoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
		cfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	}

	return cfg
}
