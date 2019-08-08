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
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gonvenience/wait"
	"github.com/homeport/gonut/internal/gonut/nok"
	"github.com/homeport/pina-golada/pkg/files"
	"github.com/mitchellh/go-homedir"
)

// PushApp performs a Cloud Foundry CLI based push operation
func PushApp(caption string, appName string, directory files.Directory, flags []string, cleanupSetting AppCleanupSetting, noPingSetting bool) (*PushReport, error) {
	if !isLoggedIn() {
		return nil, nok.Errorf(
			fmt.Sprintf("failed to push application %s to Cloud Foundry", appName),
			"session is not logged into a Cloud Foundry environment",
		)
	}

	if !isTargetOrgAndSpaceSet() {
		return nil, nok.Errorf(
			fmt.Sprintf("failed to push application %s to Cloud Foundry", appName),
			"no target is set",
		)
	}

	report := PushReport{
		AppName: appName,
	}

	err := runWithTempDir(func(path string) error {
		// Changed during each step of the verification process
		step := "Ramp-up"

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
			return nok.Errorf(
				fmt.Sprintf("failed to push application %s to Cloud Foundry", appName),
				fmt.Sprintf("An error occurred while trying to write the sample app files to disk: %v", err),
			)
		}

		pathToSampleApp := filepath.Join(path, directory.AbsolutePath().String())
		if err := os.Chdir(pathToSampleApp); err != nil {
			return nok.Errorf(
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

		// Concatenate flags and args to single slice
		args := []string{"push", appName}
		args = append(args, flags...)

		// Push application using CLI
		if output, err := cf(updates, args...); err != nil {
			caption := fmt.Sprintf("failed to push application %s to Cloud Foundry", appName)

			// Redefine caption in case Cloud Foundry gives us staging failure details
			if app, appDetailsError := getApp(appName); appDetailsError == nil {
				if app.Entity.StagingFailedDescription != nil && app.Entity.StagingFailedReason != nil {
					caption = fmt.Sprintf("%s (%s)", app.Entity.StagingFailedDescription, app.Entity.StagingFailedReason)
				}
			}

			// Try to get recent app logs to be appended to the error output
			if recentLogs, err := cf(nil, "logs", appName, "--recent"); err == nil {
				output = fmt.Sprintf("%s\n\nApplication logs:\n%s",
					output,
					recentLogs,
				)
			}

			return nok.Errorf(caption, output)
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

		// If pinging is not disabled, ping the pushed app to
		// determine its statuscode.
		if !noPingSetting {
			// Get public URL of application
			appRoute, err := getAppRoute(appName)
			if err != nil {
				return nok.Errorf(
					fmt.Sprintf("failed to get url of application %s from Cloud Foundry", appName),
					err.Error(),
				)
			}

			if statusCode, err := getAppStatusCode(appRoute); err == nil {
				if statusCode != http.StatusOK {
					return nok.Errorf(
						fmt.Sprintf("application %s returned a non-ok statuscode %d", appName, statusCode),
						"The application did not return the statuscode 200. Please try to push the same sample application again.",
					)
				}
				report.StatusCode = statusCode

			} else {
				return nok.Errorf(
					fmt.Sprintf("unable to ping application %s with route %s", appName, appRoute),
					err.Error(),
				)
			}
		}

		// If cleanup setting is set to OnSuccess, run the app removal and
		// report any issues that might come up during that operation.
		if cleanupSetting == OnSuccess {
			if output, err := cf(updates, "delete", appName, "-r", "-f"); err != nil {
				return nok.Errorf(
					fmt.Sprintf("failed to delete application %s from Cloud Foundry", appName),
					output,
				)
			}
		}

		return nil
	})

	return &report, err
}

// DeleteApps iterates over the apps in the slice and delete them
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

// HasStack returns true if Cloud Foundry that stack with the
// given name exists in the list of installed stacks
func HasStack(stackName string) (bool, error) {
	stacks, err := getStacks()
	if err != nil {
		return false, err
	}

	for _, stack := range stacks {
		if stack.Entity.Name == stackName {
			return true, nil
		}
	}

	return false, nil
}

// HasBuildpack returns true if Cloud Foundry reports that a buildpack with the
// given name exists in the list of installed buildpacks
func HasBuildpack(buildpackName string) (bool, error) {
	buildpacks, err := getBuildpacks()
	if err != nil {
		return false, err
	}

	for _, buildpack := range buildpacks {
		if buildpack.Entity.Name == buildpackName {
			return true, nil
		}
	}

	return false, nil
}

// IsExternalBuildpack returns true if the buildpackName is a valid URL which returns
// a 200 statuscode and links to a .zip file
func IsExternalBuildpack(buildpackName string) (bool, error) {
	if u, err := url.ParseRequestURI(buildpackName); err == nil {
		// Check if URL is reachable
		resp, err := http.Get(buildpackName)
		if err != nil {
			return false, err
		}
		if resp.StatusCode == http.StatusOK {
			// Check if it links to a valid buildpack file (.zip packaged)
			if strings.HasSuffix(u.Path, ".zip") {
				return true, nil
			}
			return false, fmt.Errorf("%s is not a valid buildpack! buildpacks must be .zip files", buildpackName)
		}
		return false, fmt.Errorf("could not get buildpack via URL. %s returned a non-ok statuscode", buildpackName)

	}
	return false, nil
}

func deleteApp(updates chan string, app AppDetails) error {
	if !isLoggedIn() {
		return nok.Errorf(
			fmt.Sprintf("failed to delete application %s", app.Entity.Name),
			"session is not logged into a Cloud Foundry environment",
		)
	}

	if !isTargetOrgAndSpaceSet() {
		return nok.Errorf(
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

// getAppRoute returns the public URL of the application
// using its name and the Cloud Foundry host domain.
func getAppRoute(appName string) (string, error) {
	host, domain, err := getHostAndDomain(appName)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("http://%s.%s", host, domain), nil
}

// getAppStatusCode sends a GET request to the given URL
// and returns the statuscode.
func getAppStatusCode(appRoute string) (int, error) {
	resp, err := http.Get(appRoute)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}

// GetApps gets all Apps of the targeted org and space
func GetApps() ([]AppDetails, error) {
	if !isLoggedIn() {
		return nil, nok.Errorf(
			"failed to get applications",
			"session is not logged into a Cloud Foundry environment",
		)
	}

	if !isTargetOrgAndSpaceSet() {
		return nil, nok.Errorf(
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

func getBuildpacks() ([]BuildpackDetails, error) {
	result := []BuildpackDetails{}
	nextURL := "/v2/buildpacks?results-per-page=10"

	for {
		if len(nextURL) == 0 {
			break
		}

		output, err := cf(nil, "curl", nextURL)
		if err != nil {
			return nil, err
		}

		var buildpackPage BuildpackPage
		if err := json.Unmarshal([]byte(output), &buildpackPage); err != nil {
			return nil, err
		}

		result = append(result, buildpackPage.Resources...)
		nextURL = buildpackPage.NextURL
	}

	return result, nil
}

func getStack(appName string) (*StackDetails, error) {
	app, err := getApp(appName)
	if err != nil {
		return nil, err
	}

	return cfCurlStackURL(app.Entity.StackURL)
}

func getStacks() ([]StackDetails, error) {
	result := []StackDetails{}
	nextURL := "/v2/stacks?results-per-page=10"

	for {
		if len(nextURL) == 0 {
			break
		}

		output, err := cf(nil, "curl", nextURL)
		if err != nil {
			return nil, err
		}

		var stackPage StackPage
		if err := json.Unmarshal([]byte(output), &stackPage); err != nil {
			return nil, err
		}

		result = append(result, stackPage.Resources...)
		nextURL = stackPage.NextURL
	}

	return result, nil
}

// getHostAndDomain returns the current Cloud Foundry
// hostname and domain of the application.
func getHostAndDomain(appName string) (string, string, error) {
	app, err := getApp(appName)
	if err != nil {
		return "", "", err
	}

	// Get route details of the application
	routesURL := app.Entity.RoutesURL
	routePage, err := cfCurlRouteURL(routesURL)
	if err != nil {
		return "", "", err
	}

	routeEntry := routePage.Resources[0] // Resources route always contains only one element
	host := routeEntry.Entity.Host

	// Get domain details of the application
	domainGUID := routeEntry.Entity.DomainGUID
	domainDetails, err := cfCurlDomainByGUID(domainGUID)
	if err != nil {
		return "", "", err
	}

	// Get domain string from domain details
	domain := domainDetails.Entity.Name
	return host, domain, nil
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
		line := scanner.Text()

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

func cfCurlRouteURL(routesURL string) (*RoutePage, error) {
	result, err := cf(nil, "curl", routesURL)
	if err != nil {
		return nil, err
	}

	var route RoutePage
	if err := json.Unmarshal([]byte(result), &route); err != nil {
		return nil, err
	}

	return &route, nil
}

func cfCurlDomainByGUID(domainGUID string) (*DomainDetails, error) {
	result, err := cf(nil, "curl", fmt.Sprintf("/v2/shared_domains/%s", domainGUID))
	if err != nil {
		return nil, err
	}

	var domain DomainDetails
	if err := json.Unmarshal([]byte(result), &domain); err != nil {
		return nil, err
	}

	return &domain, nil
}
