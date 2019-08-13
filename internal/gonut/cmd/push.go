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
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/gonvenience/bunt"
	"github.com/gonvenience/neat"
	"github.com/gonvenience/text"
	"github.com/homeport/gonut/internal/gonut/assets"
	"github.com/homeport/gonut/internal/gonut/cf"
	"github.com/homeport/pina-golada/pkg/files"
)

// GonutAppPrefix is the prefeix for gonuts applications, it is also used by the
// cleanup command to decide whether an app is pushed by gonut or not.
var GonutAppPrefix = "gonut"

type sampleApp struct {
	caption       string
	buildpack     string
	stack         string
	command       string
	aliases       []string
	appNamePrefix string
	assetFunc     func() (files.Directory, error)
}

var (
	deleteSetting    string
	outputSetting    string
	buildpackSetting string
	stackSetting     string
	noPingSetting    bool
)

var sampleApps = []sampleApp{
	{
		caption:       "Golang",
		command:       "golang",
		buildpack:     "go_buildpack",
		aliases:       []string{"go"},
		appNamePrefix: fmt.Sprintf("%s-golang-app-", GonutAppPrefix),
		assetFunc:     assets.Provider.GoSampleApp,
	},

	{
		caption:       "Python",
		command:       "python",
		buildpack:     "python_buildpack",
		aliases:       []string{},
		appNamePrefix: fmt.Sprintf("%s-python-app-", GonutAppPrefix),
		assetFunc:     assets.Provider.PythonSampleApp,
	},

	{
		caption:       "PHP",
		command:       "php",
		buildpack:     "php_buildpack",
		aliases:       []string{},
		appNamePrefix: fmt.Sprintf("%s-php-app-", GonutAppPrefix),
		assetFunc:     assets.Provider.PHPSampleApp,
	},

	{
		caption:       "Staticfile",
		command:       "staticfile",
		buildpack:     "staticfile_buildpack",
		aliases:       []string{"static"},
		appNamePrefix: fmt.Sprintf("%s-staticfile-app-", GonutAppPrefix),
		assetFunc:     assets.Provider.StaticfileSampleApp,
	},

	{
		caption:       "Swift",
		command:       "swift",
		buildpack:     "swift_buildpack",
		aliases:       []string{},
		appNamePrefix: fmt.Sprintf("%s-swift-app-", GonutAppPrefix),
		assetFunc:     assets.Provider.SwiftSampleApp,
	},

	{
		caption:       "NodeJS",
		command:       "nodejs",
		buildpack:     "nodejs_buildpack",
		aliases:       []string{"node"},
		appNamePrefix: fmt.Sprintf("%s-nodejs-app-", GonutAppPrefix),
		assetFunc:     assets.Provider.NodeJSSampleApp,
	},

	{
		caption:       "Ruby",
		command:       "ruby",
		buildpack:     "ruby_buildpack",
		aliases:       []string{},
		appNamePrefix: fmt.Sprintf("%s-ruby-sinatra-app-", GonutAppPrefix),
		assetFunc:     assets.Provider.RubySampleApp,
	},

	{
		caption:       ".NET",
		command:       "dotnet",
		buildpack:     "dotnet-core",
		aliases:       []string{"node"},
		appNamePrefix: fmt.Sprintf("%s-dotnet-app-", GonutAppPrefix),
		assetFunc:     assets.Provider.DotNetSampleApp,
	},

	{
		caption:       "Binary",
		command:       "binary",
		buildpack:     "binary_buildpack",
		aliases:       []string{"node"},
		appNamePrefix: fmt.Sprintf("%s-binary-app-", GonutAppPrefix),
		assetFunc:     assets.Provider.BinarySampleApp,
	},

	{
		caption:       "Java",
		command:       "java",
		buildpack:     "java_buildpack",
		aliases:       []string{"node"},
		appNamePrefix: fmt.Sprintf("%s-java-app-", GonutAppPrefix),
		assetFunc:     assets.Provider.JavaSampleApp,
	},
}

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:           "push [app]",
	Short:         "Push a sample app to Cloud Foundry",
	Long:          "Use pre-installed or remote sample apps to be pushed to a Cloud Foundry instance.",
	Example:       getOptions(),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          pushCommandFunc,
}

func init() {
	rootCmd.AddCommand(pushCmd)

	pushCmd.PersistentFlags().StringVarP(&deleteSetting, "delete", "d", "always", "Delete application after push: always, never, on-success")
	pushCmd.PersistentFlags().StringVarP(&outputSetting, "output", "o", "short", "Push output detail level: quiet, short, full")
	pushCmd.PersistentFlags().StringVarP(&buildpackSetting, "buildpack", "b", "", "Specify buildpack for pushed application")
	pushCmd.PersistentFlags().StringVarP(&stackSetting, "stack", "s", "", "Specify stack for pushed application")
	pushCmd.PersistentFlags().BoolVarP(&noPingSetting, "no-ping", "p", false, "Do not ping application after push")
}

func getOptions() string {
	options := []string{"file:<path>", "http://<git repo hostpath>", "https://<git repo hostpath>", "all"}
	for _, sampleApp := range sampleApps {
		options = append(options, sampleApp.command)
	}
	return strings.Join(options, "\n")
}

func lookUpSampleAppByName(name string) *sampleApp {
	for _, sampleApp := range sampleApps {
		if sampleApp.command == name {
			return &sampleApp
		}
	}

	return nil
}

func lookUpSampleAppByURL(absoluteURL string) *sampleApp {
	absoluteURL = strings.Trim(absoluteURL, "/") // Remove leading and tailing '/' to avoid fault paths
	rootURL := absoluteURL
	var relativePath string

	// Spilt URL into relative path and root URL (<path>#<git repo url>)
	// Example: assests/dora#github.com/cloudfoundry-samples/cf-sample-app-nodejs
	urlParts := strings.SplitN(absoluteURL, "#", 2)
	if len(urlParts) == 2 {
		relativePath = urlParts[0]
		rootURL = urlParts[1]
	}

	// Check URL validity and create sample app structure in case it is valid
	if u, err := url.ParseRequestURI(rootURL); err == nil {
		pathSlice := strings.Split(u.Path, "/")
		caption := pathSlice[len(pathSlice)-1] + "/" + relativePath

		return &sampleApp{
			caption:       caption,
			appNamePrefix: fmt.Sprintf("%s-custom-app-", GonutAppPrefix),
			assetFunc: func() (files.Directory, error) {
				// Example local path: ~/.gonut/github.com/cloudfoundry/cf-acceptance-tests/
				localPath := cf.HomeDir() + "/.gonut/" + u.Host + u.Path
				return cf.LoadSampleAppURL(rootURL, relativePath, localPath)
			},
		}
	}

	return nil
}

func pushCommandFunc(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("Please provide a sample app name as argument. Valid Arguments:\n\n%s", getOptions())
	}

	var apps []*sampleApp
	for _, arg := range args {
		if arg == "all" {
			for _, app := range sampleApps {
				apps = append(apps, &app)
			}
		} else if app := lookUpSampleAppByName(arg); app != nil {
			apps = append(apps, app)
		} else if app := lookUpSampleAppByURL(arg); app != nil {
			apps = append(apps, app)
		} else {
			return fmt.Errorf("Could not find %s sample app. Please use an argument from the following list:\n\n%s", arg, getOptions())
		}
	}

	for _, app := range apps {
		switch stackSetting {
		case "all":
			stacks, err := cf.GetStackNames()
			if err != nil {
				return fmt.Errorf("An error occurred while trying to retrieve a list of installed stacks: %v", err)
			}

			for _, stack := range stacks {
				app.stack = stack
				if err := runSampleAppPush(app); err != nil {
					return err
				}
			}
		default:
			app.stack = stackSetting // Empty if flag not set
			if err := runSampleAppPush(app); err != nil {
				return err
			}
		}

	}

	return nil
}

func runSampleAppPush(app *sampleApp) error {
	// Prepare flags for cf push command
	flags := []string{}

	// Check for stack existence
	if len(buildpackSetting) > 0 {
		app.buildpack = buildpackSetting

		hasBuildpack, err := cf.HasBuildpack(app.buildpack)
		if err != nil {
			return err
		}
		isExternalBuildpack, err := cf.IsExternalBuildpack(app.buildpack)
		if err != nil {
			return err
		}

		// Skip sample app push if desired buildpack is unavailable
		if !hasBuildpack && !isExternalBuildpack {
			bunt.Printf("Skipping push of *%s* sample app, because there is no DarkSeaGreen{%s} installed.\n",
				app.caption,
				app.buildpack,
			)
			return nil
		}

		flags = append(flags, "-b", app.buildpack)
	}

	// Check for buildpack existence for pre-defined buildpacks
	if len(app.stack) > 0 {
		hasStack, err := cf.HasStack(app.stack)
		if err != nil {
			return err
		}

		// Skip sample app push if desired stack is unavailable
		if !hasStack {
			bunt.Printf("Skipping push of *%s* sample app, because there is no DarkSeaGreen{%s} stack installed.\n",
				app.caption,
				app.buildpack,
			)
			return nil
		}

		flags = append(flags, "-s", app.stack)
	}

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

	report, err := cf.PushApp(app.caption, appName, directory, flags, cleanupSetting, noPingSetting)
	if err != nil {
		return err
	}

	switch strings.ToLower(outputSetting) {
	case "quiet":
		// Nothing to report

	case "short", "oneline":
		bunt.Printf("Successfully pushed *%s* sample app in CadetBlue{%s}.\n",
			app.caption,
			cf.HumanReadableDuration(report.ElapsedTime()),
		)

	case "json":
		out, err := neat.NewOutputProcessor(true, true, &neat.DefaultColorSchema).ToJSON(report.Export())
		if err != nil {
			return err
		}

		fmt.Println(out)

	case "yaml":
		out, err := neat.ToYAMLString(report.Export())
		if err != nil {
			return err
		}

		fmt.Println(out)

	case "full":
		headline := bunt.Sprintf("Successfully pushed *%s* sample app in CadetBlue{%s}",
			app.caption,
			cf.HumanReadableDuration(report.ElapsedTime()),
		)

		content, err := neat.Table(report.ExportTable(), neat.AlignRight(0))
		if err != nil {
			return err
		}

		neat.Box(os.Stdout, headline, strings.NewReader(content))
	}

	return nil
}
