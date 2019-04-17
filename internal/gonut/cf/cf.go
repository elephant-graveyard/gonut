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

package cf

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/homeport/gonvenience/pkg/v1/bunt"
	"github.com/homeport/gonvenience/pkg/v1/wait"
	"github.com/homeport/pina-golada/pkg/files"
	"github.com/mitchellh/go-homedir"
)

// PushApp performs a Cloud Foundry CLI based push operation
func PushApp(caption string, appName string, directory files.Directory, cleanupSetting AppCleanupSetting) (*PushReport, error) {
	if !isLoggedIn() {
		return nil, Errorf(
			fmt.Sprintf("failed to push application %s to Cloud Foundry", appName),
			"session is not logged into a Cloud Foundry environment",
		)
	}

	if !isTargetOrgAndSpaceSet() {
		return nil, Errorf(
			fmt.Sprintf("failed to push application %s to Cloud Foundry", appName),
			"no target is set",
		)
	}

	report := PushReport{}

	err := runWithTempDir(func(path string) error {
		// Changed during each step of the verification process
		step := "Initialisation"

		spinner := wait.NewProgressIndicator("*%s*, DimGray{%s}", caption, step)
		spinner.Start()
		defer spinner.Stop()

		updates := make(chan string)
		defer close(updates)
		go func() {
			for update := range updates {
				if text := strings.Trim(update, " "); len(text) > 0 {
					if result := report.ParseUpdate(text); result != "" {
						step = result
					}

					spinner.SetText("*%s*, DimGray{%s} - %s",
						caption,
						step,
						text,
					)
				}
			}
		}()

		if err := files.WriteToDisk(directory, path, true); err != nil {
			return Errorf(
				fmt.Sprintf("failed to push application %s to Cloud Foundry", appName),
				fmt.Sprintf("An error occurred while trying to write the sample app files to disk: %v", err),
			)
		}

		pathToSampleApp := filepath.Join(path, directory.AbsolutePath().String())
		if err := os.Chdir(pathToSampleApp); err != nil {
			return Errorf(
				fmt.Sprintf("failed to push application %s to Cloud Foundry", appName),
				fmt.Sprintf("An error occurred while trying to change working directory to %s: %v", pathToSampleApp, err),
			)
		}

		// If cleanup setting is set to always, make sure to run the delete app
		// CF CLI call no matter what happens next.
		if cleanupSetting == Always {
			defer cf(updates, "delete", appName, "-r", "-f")
		}

		// Note the timestamp when the push starts
		report.InitStart = time.Now()

		if output, err := cf(updates, "push", appName); err != nil {
			if app, appDetailsError := getApp(appName); appDetailsError == nil {
				if app.Entity.StagingFailedDescription != nil && app.Entity.StagingFailedReason != nil {
					return Errorf(
						fmt.Sprintf("failed to push application %s to Cloud Foundry", appName),
						fmt.Sprintf("%s (%s)", app.Entity.StagingFailedDescription, app.Entity.StagingFailedReason),
					)
				}
			}

			if recentLogs, err := cf(nil, "logs", appName, "--recent"); err == nil {
				output = fmt.Sprintf("%s\n\nApplication logs:\n%s",
					output,
					recentLogs,
				)
			}

			return Errorf(
				fmt.Sprintf("failed to push application %s to Cloud Foundry", appName),
				output,
			)
		}

		// Note the timestamp when the push has finished
		report.PushEnd = time.Now()

		// Gather details about the buildpack used for the app
		if buildpack, err := getBuildpack(appName); err == nil {
			report.buildpack = buildpack
		}

		// Gather details about the stack used for the app
		if stack, err := getStack(appName); err == nil {
			report.stack = stack
		}

		// If cleanup setting is set to OnSuccess, run the app removal and
		// report any issues that might come up during that operation.
		if cleanupSetting == OnSuccess {
			if output, err := cf(updates, "delete", appName, "-r", "-f"); err != nil {
				return Errorf(
					fmt.Sprintf("failed to delete application %s from Cloud Foundry", appName),
					output,
				)
			}
		}

		return nil
	})

	return &report, err
}

//DeleteApps iterates over the apps in the slice and delete them
func DeleteApps(apps []AppDetails) error {
	caption := "Cleaning Up"
	spinner := wait.NewProgressIndicator("*%s*", caption)
	spinner.Start()
	defer spinner.Stop()

	updates := make(chan string)
	defer close(updates)
	go func() {
		for update := range updates {
			if text := strings.Trim(update, " "); len(text) > 0 {
				spinner.SetText("*%s*, %s",
					caption,
					text,
				)
			}
		}
	}()

	for _, app := range apps {
		if err := deleteApp(updates, app); err != nil {
			return err
		}
	}
	return nil
}

func deleteApp(updates chan string, app AppDetails) error {

	if !isLoggedIn() {
		return Errorf(
			fmt.Sprintf("failed to delete application %s", app.Entity.Name),
			"session is not logged into a Cloud Foundry environment",
		)
	}

	if !isTargetOrgAndSpaceSet() {
		return Errorf(
			fmt.Sprintf("failed to delete application %s", app.Entity.Name),
			"no target is set",
		)
	}

	if _, err := cf(updates, "delete", app.Entity.Name, "-r", "-f"); err != nil {
		return err
	}

	return nil
}

func runWithTempDir(f func(path string) error) error {
	dir, err := ioutil.TempDir("", "gonut")
	if err != nil {
		return err
	}

	defer os.RemoveAll(dir)

	return f(dir)
}

func isLoggedIn() bool {
	config, err := getCloudFoundryConfig()
	if err != nil {
		return false
	}

	return len(config.AccessToken) > 0
}

func isTargetOrgAndSpaceSet() bool {
	org, space, err := getOrgAndSpaceNamesFromConfig()
	return len(org) > 0 && len(space) > 0 && err == nil
}

func getOrgAndSpaceNamesFromConfig() (string, string, error) {
	config, err := getCloudFoundryConfig()
	if err != nil {
		return "", "", err
	}

	return config.OrganizationFields.Name, config.SpaceFields.Name, nil
}

func getCloudFoundryConfig() (*CloudFoundryConfig, error) {
	path, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(filepath.Join(path, ".cf", "config.json"))
	if err != nil {
		return nil, err
	}

	var config CloudFoundryConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func getApp(appName string) (*AppDetails, error) {
	appGUID, err := cfAppGUID(appName)
	if err != nil {
		return nil, err
	}

	return cfCurlAppByGUID(appGUID)
}

//GetApps gets all Apps of the targeted org and space
func GetApps() ([]AppDetails, error) {
	if !isLoggedIn() {
		return nil, Errorf(
			"failed to get applications",
			"session is not logged into a Cloud Foundry environment",
		)
	}

	if !isTargetOrgAndSpaceSet() {
		return nil, Errorf(
			"failed to get applications",
			"no target is set",
		)
	}

	return cfCurlOrgsSpaceApps()
}

func getBuildpack(appName string) (*BuildpackDetails, error) {
	app, err := getApp(appName)
	if err != nil {
		return nil, err
	}

	return cfCurlBuildpackByGUID(app.Entity.DetectedBuildpackGUID)
}

func getStack(appName string) (*StackDetails, error) {
	app, err := getApp(appName)
	if err != nil {
		return nil, err
	}

	return cfCurlStackURL(app.Entity.StackURL)
}

func cf(updates chan string, args ...string) (string, error) {
	var (
		buf bytes.Buffer
		err error
	)

	read, write := io.Pipe()
	go func() {
		// TODO Check if cf binary is available
		cmd := exec.Command("cf", args...)

		cmd.Stdout = write
		cmd.Stderr = write

		err = cmd.Run()
		write.Close()
	}()

	scanner := bufio.NewScanner(read)

	for scanner.Scan() {
		// TODO As long as the bunt package cannot handle arbitrary long strings
		// with color markers, remove the color to avoid panics during runtime.
		line := bunt.RemoveAllEscapeSequences(scanner.Text())

		// Some cloud controller versions return unwanted control characters
		// that are better be removed to achieve a clean output
		line = strings.Replace(line, "\r", "", -1)

		buf.WriteString(line)
		buf.WriteString("\n")

		if updates != nil {
			updates <- line
		}
	}

	return buf.String(), err
}

func cfAppGUID(appName string) (string, error) {
	result, err := cf(nil, "app", appName, "--guid")
	if err != nil {
		return "", err
	}

	return strings.Trim(result, " \n"), nil
}

func cfCurlOrgsSpaceApps() ([]AppDetails, error) {
	var apps []AppDetails
	var result string
	var err error

	nextURL := "/v2/apps"

	for {
		result, err = cf(nil, "curl", nextURL)

		if err != nil {
			return nil, err
		}

		var page AppsPage

		if err := json.Unmarshal([]byte(result), &page); err != nil {
			return nil, err
		}

		apps = append(apps, page.Resources...)

		if len(page.NextURL) == 0 {
			break
		}

		nextURL = page.NextURL
	}

	return apps, nil
}

func cfCurlAppByGUID(appGUID string) (*AppDetails, error) {
	result, err := cf(nil, "curl", fmt.Sprintf("/v2/apps/%s", appGUID))
	if err != nil {
		return nil, err
	}

	var app AppDetails
	if err := json.Unmarshal([]byte(result), &app); err != nil {
		return nil, err
	}

	return &app, nil
}

func cfCurlBuildpackByGUID(buildpackGUID string) (*BuildpackDetails, error) {
	result, err := cf(nil, "curl", fmt.Sprintf("/v2/buildpacks/%s", buildpackGUID))
	if err != nil {
		return nil, err
	}

	var buildpack BuildpackDetails
	if err := json.Unmarshal([]byte(result), &buildpack); err != nil {
		return nil, err
	}

	return &buildpack, nil
}

func cfCurlStackURL(stackURL string) (*StackDetails, error) {
	result, err := cf(nil, "curl", stackURL)
	if err != nil {
		return nil, err
	}

	var stack StackDetails
	if err := json.Unmarshal([]byte(result), &stack); err != nil {
		return nil, err
	}

	return &stack, nil
}
