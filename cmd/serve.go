package cmd

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/leotaku/maiden/caldav"
	"github.com/leotaku/maiden/diary"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:     "serve [flags..]",
	Short:   "Start task synchronization daemon",
	Version: rootCmd.Version,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		order, err := validateDateStyle(dateStyleArg)
		if err != nil {
			return fmt.Errorf("args: %w", err)
		}

		diary, err := diary.NewDiary(fileArg, order)
		if err != nil {
			return fmt.Errorf("diary: %w", err)
		}

		client, err := loadGoogleClient(calendarArg, dataHome)
		if err != nil {
			return fmt.Errorf("caldav: %w", err)
		}

		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			return fmt.Errorf("watch: %w", err)
		}
		defer watcher.Close()

		if err := watcher.Add(fileArg); err != nil {
			return fmt.Errorf("watch %v: %w", fileArg, err)
		}

		// All errors are only logged past this point
		logrus.New().Info("Synchronization daemon started...")
		sync := client.Sync(time.Second)
		for {
			ctag, _ := client.GetCtag()
			select {
			case <-sync.Wait(ctag):
				logrus.New().Info("Calendar changed, starting synchronization")
				updateAndSync(client, diary)
				logrus.New().Info("Synchronization finished")
			case <-watchNow(watcher):
				logrus.New().Info("Diary changed, starting synchronization")
				updateAndSync(client, diary)
				logrus.New().Info("Synchronization finished")
			}
		}

	},
	DisableFlagsInUseLine: true,
}

func updateAndSync(c *caldav.Client, d *diary.Diary) {
	if err := c.Update(); err != nil {
		return
	}
	if err := d.Update(); err != nil {
		return
	}

	syncOnce(c, d)
}

func watchNow(watcher *fsnotify.Watcher) <-chan fsnotify.Event {
	output := make(chan fsnotify.Event)
	go func() {
		for {
			select {
			case <-watcher.Events:
			case <-watcher.Errors:
			default:
				output <- <-watcher.Events
				close(output)
				return
			}
		}
	}()

	return output
}

func initServeCmd() {
	*serveCmd.Flags() = *syncCmd.Flags()
}
