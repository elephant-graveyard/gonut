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

var sampleApps = []sampleApp{
	sampleApp{
		caption:       "Golang",
		command:       "golang",
		aliases:       []string{"go"},
		appNamePrefix: "gonut-golang-app-",
		assetFunc:     assets.Provider.GoSampleApp,
	},

	sampleApp{
		caption:       "Python",
		command:       "python",
		aliases:       []string{},
		appNamePrefix: "gonut-python-app-",
		assetFunc:     assets.Provider.PythonSampleApp,
	},

	sampleApp{
		caption:       "PHP",
		command:       "php",
		aliases:       []string{},
		appNamePrefix: "gonut-php-app-",
		assetFunc:     assets.Provider.PHPSampleApp,
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

	for _, sampleApp := range sampleApps {
		pushCmd.AddCommand(&cobra.Command{
			Use:   sampleApp.command,
			Short: fmt.Sprintf("Push a %s sample app to Cloud Foundry", sampleApp.caption),
			Long:  fmt.Sprintf(`Push a %s sample app to Cloud Foundry. The application will be deleted after it was pushed successfully.`, sampleApp.caption),
			RunE:  genericCommandFunc,
		})
	}
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

	appName := text.RandomStringWithPrefix(sampleApp.appNamePrefix, 32)

	directory, err := sampleApp.assetFunc()
	if err != nil {
		return err
	}

	return cf.PushApp(sampleApp.caption, appName, directory)
}
