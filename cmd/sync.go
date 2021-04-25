package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/leotaku/maiden/caldav"
	"github.com/leotaku/maiden/diary"
	"github.com/leotaku/maiden/oauth"
	"github.com/spf13/cobra"
)

var (
	calendarArg  string
	fileArg      string
	dateStyleArg string
)

var serveCmd = &cobra.Command{
	Use:     "serve [flags..]",
	Short:   "Start task synchronization daemon",
	Version: rootCmd.Version,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		life := oauth.NewLifecycle(dataHome)
		config, tok, err := life.Load()
		if err != nil {
			return fmt.Errorf("auth: %w", err)
		}

		href := fmt.Sprintf("/caldav/v2/%v/events", calendarArg)
		http := config.Client(context.TODO(), tok)
		b := caldav.NewBuilder()
		b.WithHttp(http)
		b.WithHostURL(googleURL)
		b.WithCalendarPath(href)
		client, err := b.BuildAndInit()
		if err != nil {
			return fmt.Errorf("init: %w", err)
		}

		for _, event := range client.Events() {
			if strings.Contains(event.Id(), client.ProviderName()) {
				fmt.Println(client.Del(event))
			}
		}

		f, err := os.Open(fileArg)
		if err != nil {
			return fmt.Errorf("local: %w", err)
		}
		entries, err := diary.NewParser(f, diary.ISO).All()
		if err != nil {
			return fmt.Errorf("diary: %w", err)
		}

		for _, entry := range entries {
			fmt.Printf("%-20v %-35v %v\n", entry.Description, entry.Timestamp.First(time.Local), entry.Duration)
			fmt.Println(client.Put(entry.ToICALEvent(time.Local)))
		}

		return nil
	},
	DisableFlagsInUseLine: true,
}

func init() {
	serveCmd.Flags().StringVarP(&calendarArg, "calendar", "c", "", "Calendar ID to connect to")
	serveCmd.Flags().StringVarP(&fileArg, "file", "f", "", "Local diary file to use")
	serveCmd.Flags().StringVarP(&dateStyleArg, "date-style", "d", "", "Date style used in diary")
	serveCmd.MarkFlagRequired("calendar")
	serveCmd.MarkFlagRequired("file")
	serveCmd.MarkFlagRequired("date-style")
}
