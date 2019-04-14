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
	Run:   cleanUp,
}

func init() {
	rootCmd.AddCommand(cleanUpCmd)
}

func cleanUp(cmd *cobra.Command, args []string) {
	apps, _ := cf.GetApps()
	var appsToClean, _ = getGonutApps(apps, GonutAppPrefix)
	if err := cf.DeleteApps(appsToClean); err != nil {

	}
}

func getGonutApps(apps []cf.AppDetails, prefix string) ([]cf.AppDetails, error) {
	gonutApps := apps[:0]
	for _, app := range apps {

		if strings.HasPrefix(app.Entity.Name, prefix) {
			gonutApps = append(gonutApps, app)
		}
	}
	return gonutApps, nil
}
