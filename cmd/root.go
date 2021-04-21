package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	noColorArg bool
)

var rootCmd = &cobra.Command{
	Use:                   "maiden [flags..] <subcommand>",
	Short:                 "The Maiden task synchronization system",
	Version:               "0.1",
	DisableFlagsInUseLine: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(initCmd, serveCmd)
	rootCmd.SetHelpFunc(help)
	rootCmd.SetUsageFunc(usage)
	_, noColor := os.LookupEnv("NO_COLOR")
	rootCmd.Flags().BoolVarP(&noColorArg, "no-color", "n", noColor, "whether to disable colors in output")
	rootCmd.Flags().SortFlags = false
}
