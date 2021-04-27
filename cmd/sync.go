package cmd

import (
	"fmt"
	"strings"

	"github.com/leotaku/maiden/caldav"
	"github.com/leotaku/maiden/diary"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	calendarArg  string
	fileArg      string
	dateStyleArg string
)

var syncCmd = &cobra.Command{
	Use:     "sync [flags..]",
	Short:   "Sync tasks exactly once",
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

		// All errors are only logged past this point
		syncOnce(client, diary)

		return nil
	},
	DisableFlagsInUseLine: true,
}

func syncOnce(c *caldav.Client, d *diary.Diary) {
	for _, event := range c.Events() {
		if !strings.Contains(event.Id(), c.ProviderName()) {
			if entry, err := diary.FromVEvent(event); err != nil {
				logrus.WithFields(logrus.Fields{
					"id":    event.Id(),
					"error": err,
				}).Warn("Importing event failed")
				continue
			} else if err := d.Add(*entry); err != nil {
				logrus.WithFields(logrus.Fields{
					"id":    event.Id(),
					"error": err,
				}).Warn("Writing event to diary failed")
				continue
			}
		}
		if err := c.Del(event); err != nil {
			logrus.WithFields(logrus.Fields{
				"id":    event.Id(),
				"error": err,
			}).Error("Deletion of event failed")
		}
	}

	_, loc := c.Timezone()
	for _, entry := range d.Entries() {
		ctx := logrus.WithFields(logrus.Fields{
			"description": entry.Description,
			"timestamp":   entry.Datetime.First(loc),
			"duration":    entry.Duration,
		})
		err := c.Put(entry.ToVEvent(loc))
		if err != nil {
			ctx.WithField("error", err).Error("Upload of event failed")
		} else {
			ctx.Debug("Upload of event succeeded")
		}
	}

}

func init() {
	syncCmd.Flags().StringVarP(&calendarArg, "calendar", "c", "", "Calendar ID to connect to")
	syncCmd.Flags().StringVarP(&fileArg, "file", "f", "", "Local diary file to use")
	syncCmd.Flags().StringVarP(&dateStyleArg, "date-style", "d", "", "Date style used in diary")
	syncCmd.MarkFlagRequired("calendar")   //nolint:errcheck
	syncCmd.MarkFlagRequired("file")       //nolint:errcheck
	syncCmd.MarkFlagRequired("date-style") //nolint:errcheck
	syncCmd.Flags().SortFlags = false
	initServeCmd()
}
