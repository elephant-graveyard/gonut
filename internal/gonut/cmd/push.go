// Copyright © 2019 The Homeport Team
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/homeport/gonut/internal/gonut/assets"
	"github.com/homeport/gonut/internal/gonut/cf"
	"github.com/homeport/gonvenience/pkg/v1/text"
	"github.com/homeport/pina-golada/pkg/files"
)

type sampleApp struct {
	caption       string
	command       string
	aliases       []string
	appNamePrefix string
	assetFunc     func() (files.Directory, error)
}

var deleteSetting string

var sampleApps = []sampleApp{
	{
		caption:       "Golang",
		command:       "golang",
		aliases:       []string{"go"},
		appNamePrefix: "gonut-golang-app-",
		assetFunc:     assets.Provider.GoSampleApp,
	},

	{
		caption:       "Python",
		command:       "python",
		aliases:       []string{},
		appNamePrefix: "gonut-python-app-",
		assetFunc:     assets.Provider.PythonSampleApp,
	},

	{
		caption:       "PHP",
		command:       "php",
		aliases:       []string{},
		appNamePrefix: "gonut-php-app-",
		assetFunc:     assets.Provider.PHPSampleApp,
	},

	{
		caption:       "Staticfile",
		command:       "staticfile",
		aliases:       []string{"static"},
		appNamePrefix: "gonut-staticfile-app-",
		assetFunc:     assets.Provider.StaticfileSampleApp,
	},

	{
		caption:       "Swift",
		command:       "swift",
		aliases:       []string{},
		appNamePrefix: "gonut-swift-app-",
		assetFunc:     assets.Provider.SwiftSampleApp,
	},

	{
		caption:       "NodeJS",
		command:       "nodejs",
		aliases:       []string{"node"},
		appNamePrefix: "gonut-nodejs-app-",
		assetFunc:     assets.Provider.NodeJSSampleApp,
	},
}

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push a sample app to Cloud Foundry",
	Long:  `Use one of the sub-commands to select a sample app of a list of programming languages to be pushed to a Cloud Foundry instance.`,
}

func init() {
	rootCmd.AddCommand(pushCmd)

	pushCmd.PersistentFlags().StringVarP(&deleteSetting, "delete", "d", "always", "Delete application after push: always, never, on-success")

	for _, sampleApp := range sampleApps {
		pushCmd.AddCommand(&cobra.Command{
			Use:     sampleApp.command,
			Aliases: sampleApp.aliases,
			Short:   fmt.Sprintf("Push a %s sample app to Cloud Foundry", sampleApp.caption),
			Long:    fmt.Sprintf(`Push a %s sample app to Cloud Foundry. The application will be deleted after it was pushed successfully.`, sampleApp.caption),
			RunE:    genericCommandFunc,
		})
	}

	pushCmd.AddCommand(&cobra.Command{
		Use:   "all",
		Short: "Pushes all available sample apps to Cloud Foundry",
		Long:  `Pushes all available sample apps to Cloud Foundry. Each application will be deleted after it was pushed successfully.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, sampleApp := range sampleApps {
				if err := runSampleAppPush(sampleApp); err != nil {
					return err
				}
			}

			return nil
		},
	})
}

func lookUpSampleAppByName(name string) *sampleApp {
	for _, sampleApp := range sampleApps {
		if sampleApp.command == name {
			return &sampleApp
		}
	}

	return nil
}

func genericCommandFunc(cmd *cobra.Command, args []string) error {
	sampleApp := lookUpSampleAppByName(cmd.Use)
	if sampleApp == nil {
		return fmt.Errorf("failed to detect which sample app is to be tested")
	}

	return runSampleAppPush(*sampleApp)
}

func runSampleAppPush(app sampleApp) error {
	var cleanupSetting cf.AppCleanupSetting
	switch deleteSetting {
	case "always":
		cleanupSetting = cf.Always

	case "never":
		cleanupSetting = cf.Never

	case "on-success":
		cleanupSetting = cf.OnSuccess

	default:
		return fmt.Errorf("unsupported delete setting: %s", deleteSetting)
	}

	appName := text.RandomStringWithPrefix(app.appNamePrefix, 32)

	directory, err := app.assetFunc()
	if err != nil {
		return err
	}

	return cf.PushApp(app.caption, appName, directory, cleanupSetting)
}
