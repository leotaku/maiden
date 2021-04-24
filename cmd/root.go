package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	noColorArg  bool
	logLevelArg string
)

var rootCmd = &cobra.Command{
	Use:     "maiden [flags..] <subcommand>",
	Short:   "The Maiden task synchronization system",
	Version: "0.1",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("Subcommand required")
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logrus.SetFormatter(&logrus.TextFormatter{DisableColors: noColorArg})
		if lvl, err := logrus.ParseLevel(logLevelArg); err == nil {
			logrus.SetLevel(lvl)
		}
	},
	DisableFlagsInUseLine: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(initCmd, syncCmd, serveCmd)
	rootCmd.SetHelpFunc(help)
	rootCmd.SetUsageFunc(usage)
	_, noColor := os.LookupEnv("NO_COLOR")
	rootCmd.PersistentFlags().BoolVarP(&noColorArg, "no-color", "n", noColor, "whether to disable colors in output")
	rootCmd.PersistentFlags().StringVarP(&logLevelArg, "log", "l", "warning", "default logging level")
	rootCmd.PersistentFlags().SortFlags = false
	rootCmd.Flags().SortFlags = false
}
