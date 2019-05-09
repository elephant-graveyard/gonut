package cmd

import (
	"strings"

	"github.com/homeport/gonut/internal/gonut/cf"
	"github.com/spf13/cobra"
)

// cleanUpCmd represents the cleanup command
var cleanUpCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Delete all gonut pushed apps",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cleanUp(cmd, args); err != nil {
			ExitGonut(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(cleanUpCmd)
}

func cleanUp(cmd *cobra.Command, args []string) error {
	apps, err := cf.GetApps()
	if err != nil {
		return err
	}
	appsToClean := getGonutApps(apps, GonutAppPrefix)
	if len(appsToClean) == 0 {
		return nil
	}

	if err := cf.DeleteApps(appsToClean); err != nil {
		return err
	}

	return nil
}

func getGonutApps(apps []cf.AppDetails, prefix string) []cf.AppDetails {
	gonutApps := apps[:0]
	for _, app := range apps {

		if strings.HasPrefix(app.Entity.Name, prefix) {
			gonutApps = append(gonutApps, app)
		}
	}
	return gonutApps
}
