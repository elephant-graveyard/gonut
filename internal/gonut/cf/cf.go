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
	"github.com/pkg/errors"
)

// PushApp performs a Cloud Foundry CLI based push operation
func PushApp(caption string, appName string, directory files.Directory, cleanupSetting AppCleanupSetting) (*PushReport, error) {
	if !isLoggedIn() {
		return nil, fmt.Errorf("unable to push, session is not logged into a Cloud Foundry environment")
	}

	if !isTargetOrgAndSpaceSet() {
		return nil, fmt.Errorf("unable to push, no target is set")
	}

	report := PushReport{}

	err := runWithTempDir(func(path string) error {
		report.InitStart = time.Now()

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
					switch {
					case strings.HasPrefix(text, "Creating app"):
						report.CreatingStart = time.Now()
						step = "Creating"

					case strings.HasPrefix(text, "Uploading files"):
						report.UploadingStart = time.Now()
						step = "Uploading"

					case strings.HasPrefix(text, "Staging app"):
						report.StagingStart = time.Now()
						step = "Staging"

					case strings.HasPrefix(text, "Waiting for app to start..."):
						report.StartingStart = time.Now()
						step = "Starting"

					case strings.HasPrefix(text, "Deleting app"):
						step = "Deleting"
					}

					spinner.SetText("*%s*, DimGray{%s} - %s",
						caption,
						step,
						text,
					)
				}
			}
		}()

		pathToSampleApp := filepath.Join(path, directory.AbsolutePath().String())

		if err := files.WriteToDisk(directory, path, true); err != nil {
			return errors.Wrap(err, "failed to write sample app files to disk")
		}

		if err := os.Chdir(pathToSampleApp); err != nil {
			return errors.Wrapf(err, "failed to change working directory to %s", pathToSampleApp)
		}

		// If cleanup setting is set to always, make sure to run the delete app
		// CF CLI call no matter what happens next.
		if cleanupSetting == Always {
			defer cf(updates, "delete", appName, "-r", "-f")
		}

		if _, err := cf(updates, "push", appName); err != nil {
			return errors.Wrapf(err, "failed to push application %s to Cloud Foundry", appName)
		}

		// Note the timestamp when the push has finished
		report.PushEnd = time.Now()

		// Gather details about the buildpack used for the app
		if buildpack, err := getBuildpack(appName); err == nil {
			report.Buildpack = buildpack
		}

		// If cleanup setting is set to OnSuccess, run the app removal and
		// report any issues that might come up during that operation.
		if cleanupSetting == OnSuccess {
			if _, err := cf(updates, "delete", appName, "-r", "-f"); err != nil {
				return errors.Wrapf(err, "failed to delete application %s from Cloud Foundry", appName)
			}
		}

		return nil
	})

	return &report, err
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

func getBuildpack(appName string) (*BuildpackDetails, error) {
	appGUID, err := cfAppGUID(appName)
	if err != nil {
		return nil, err
	}

	app, err := cfCurlAppByGUID(appGUID)
	if err != nil {
		return nil, err
	}

	return cfCurlBuildpackByGUID(app.Entity.DetectedBuildpackGUID)
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
		buf.WriteString(line)

		if updates != nil {
			updates <- line
		}
	}

	return buf.String(), err
}

func cfAppGUID(appName string) (string, error) {
	return cf(nil, "app", appName, "--guid")
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
