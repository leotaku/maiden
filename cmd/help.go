package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func writeHelp(cmd *cobra.Command, w io.Writer) {
	fmt.Fprintf(w, "Usage:\n  %v\n", cmd.Use)

	if cmd.HasSubCommands() {
		fmt.Fprintf(w, "\nSubcommands:\n")
		for _, sub := range cmd.Commands() {
			fmt.Fprintf(w, "  %-20v%v\n", sub.Name(), toSentenceCase(sub.Short))
		}
	}

	fmt.Fprintf(w, "\nOptions:\n")
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		shorthand := ""
		if len(f.Shorthand) > 0 {
			shorthand = "-" + f.Shorthand + ", "
		}
		fmt.Fprintf(w, "  %4v--%-14v%v\n", shorthand, f.Name, toSentenceCase(f.Usage))
	})
}

func help(cmd *cobra.Command, args []string) {
	fmt.Fprintf(os.Stdout, "%v\n", cmd.Short)
	writeHelp(cmd, os.Stdout)
}

func usage(cmd *cobra.Command) error {
	writeHelp(cmd, os.Stderr)
	return nil
}

func toSentenceCase(sentence string) string {
	words := strings.Split(sentence, " ")
	words[0] = strings.Title(words[0])
	return strings.Join(words, " ")
}
