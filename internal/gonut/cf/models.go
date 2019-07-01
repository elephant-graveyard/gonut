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
	"time"
)

// AppCleanupSetting specifies supported cleanup options
type AppCleanupSetting int

// Supported cleanup settings include:
// - Never, leaves the app no matter if push worked or not
// - Always, removes the app after a push attempt
// - OnSuccess, removes the app if the push went through without issues
const (
	Never = AppCleanupSetting(iota)
	Always
	OnSuccess
)

// CloudFoundryConfig defines the structure used by the Cloud Foundry CLI configuration JSONs
type CloudFoundryConfig struct {
	ConfigVersion         int    `json:"ConfigVersion"`
	Target                string `json:"Target"`
	APIVersion            string `json:"APIVersion"`
	AuthorizationEndpoint string `json:"AuthorizationEndpoint"`
	DopplerEndPoint       string `json:"DopplerEndPoint"`
	UaaEndpoint           string `json:"UaaEndpoint"`
	RoutingAPIEndpoint    string `json:"RoutingAPIEndpoint"`
	AccessToken           string `json:"AccessToken"`
	SSHOAuthClient        string `json:"SSHOAuthClient"`
	UAAOAuthClient        string `json:"UAAOAuthClient"`
	UAAOAuthClientSecret  string `json:"UAAOAuthClientSecret"`
	RefreshToken          string `json:"RefreshToken"`
	OrganizationFields    struct {
		GUID            string `json:"GUID"`
		Name            string `json:"Name"`
		QuotaDefinition struct {
			GUID                    string `json:"guid"`
			Name                    string `json:"name"`
			MemoryLimit             int    `json:"memory_limit"`
			InstanceMemoryLimit     int    `json:"instance_memory_limit"`
			TotalRoutes             int    `json:"total_routes"`
			TotalServices           int    `json:"total_services"`
			NonBasicServicesAllowed bool   `json:"non_basic_services_allowed"`
			AppInstanceLimit        int    `json:"app_instance_limit"`
			TotalReservedRoutePorts int    `json:"total_reserved_route_ports"`
		} `json:"QuotaDefinition"`
	} `json:"OrganizationFields"`
	SpaceFields struct {
		GUID     string `json:"GUID"`
		Name     string `json:"Name"`
		AllowSSH bool   `json:"AllowSSH"`
	} `json:"SpaceFields"`
	SSLDisabled  bool   `json:"SSLDisabled"`
	AsyncTimeout int    `json:"AsyncTimeout"`
	Trace        string `json:"Trace"`
	ColorEnabled string `json:"ColorEnabled"`
	Locale       string `json:"Locale"`
	PluginRepos  []struct {
		Name string `json:"Name"`
		URL  string `json:"URL"`
	} `json:"PluginRepos"`
	MinCLIVersion            string `json:"MinCLIVersion"`
	MinRecommendedCLIVersion string `json:"MinRecommendedCLIVersion"`
}

// AppDetails is the Go struct for the /v2/apps/<guid> result JSON
type AppDetails struct {
	Metadata struct {
		GUID      string    `json:"guid"`
		URL       string    `json:"url"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"metadata"`
	Entity struct {
		Name                     string      `json:"name"`
		Production               bool        `json:"production"`
		SpaceGUID                string      `json:"space_guid"`
		StackGUID                string      `json:"stack_guid"`
		Buildpack                interface{} `json:"buildpack"`
		DetectedBuildpack        string      `json:"detected_buildpack"`
		DetectedBuildpackGUID    string      `json:"detected_buildpack_guid"`
		EnvironmentJSON          interface{} `json:"environment_json"`
		Memory                   int         `json:"memory"`
		Instances                int         `json:"instances"`
		DiskQuota                int         `json:"disk_quota"`
		State                    string      `json:"state"`
		Version                  string      `json:"version"`
		Command                  interface{} `json:"command"`
		Console                  bool        `json:"console"`
		Debug                    interface{} `json:"debug"`
		StagingTaskID            string      `json:"staging_task_id"`
		PackageState             string      `json:"package_state"`
		HealthCheckType          string      `json:"health_check_type"`
		HealthCheckTimeout       interface{} `json:"health_check_timeout"`
		HealthCheckHTTPEndpoint  interface{} `json:"health_check_http_endpoint"`
		StagingFailedReason      interface{} `json:"staging_failed_reason"`
		StagingFailedDescription interface{} `json:"staging_failed_description"`
		Diego                    bool        `json:"diego"`
		DockerImage              interface{} `json:"docker_image"`
		DockerCredentials        struct {
			Username interface{} `json:"username"`
			Password interface{} `json:"password"`
		} `json:"docker_credentials"`
		PackageUpdatedAt     time.Time `json:"package_updated_at"`
		DetectedStartCommand string    `json:"detected_start_command"`
		EnableSSH            bool      `json:"enable_ssh"`
		Ports                []int     `json:"ports"`
		SpaceURL             string    `json:"space_url"`
		StackURL             string    `json:"stack_url"`
		RoutesURL            string    `json:"routes_url"`
		EventsURL            string    `json:"events_url"`
		ServiceBindingsURL   string    `json:"service_bindings_url"`
		RouteMappingsURL     string    `json:"route_mappings_url"`
	} `json:"entity"`
}

// BuildpackDetails is the Go struct for the /v2/buildpacks/<guid> result JSON
type BuildpackDetails struct {
	Metadata struct {
		GUID      string    `json:"guid"`
		URL       string    `json:"url"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"metadata"`
	Entity struct {
		Name     string      `json:"name"`
		Stack    interface{} `json:"stack"`
		Position int         `json:"position"`
		Enabled  bool        `json:"enabled"`
		Locked   bool        `json:"locked"`
		Filename string      `json:"filename"`
	} `json:"entity"`
}

// StackDetails is the Go struct for the /v2/stacks/<guid> result JSON
type StackDetails struct {
	Metadata struct {
		GUID      string    `json:"guid"`
		URL       string    `json:"url"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"metadata"`
	Entity struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"entity"`
}

// AppsPage represents the result from cf curl /v2/apps page
type AppsPage struct {
	// TotalResults string `json:"total_results"`
	// TotalPages string `json:"total_pages"`
	// PrevURL     string       `json:"prev_url"`
	NextURL   string       `json:"next_url"`
	Resources []AppDetails `json:"resources"`
}

// BuildpackPage represents the result of cf curl /v2/buildpacks output
type BuildpackPage struct {
	TotalResults int                `json:"total_results"`
	TotalPages   int                `json:"total_pages"`
	PrevURL      string             `json:"prev_url"`
	NextURL      string             `json:"next_url"`
	Resources    []BuildpackDetails `json:"resources"`
}

// RouteDetails is the Go struct for the /v2/apps/<guid>/routes result JSON
type RouteDetails struct {
	Metadata struct {
		GUID      string    `json:"guid"`
		URL       string    `json:"url"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"metadata"`
	Entity struct {
		Host                string      `json:"host"`
		Path                string      `json:"path"`
		DomainGUID          string      `json:"domain_guid"`
		SpaceGUID           string      `json:"space_guid"`
		ServiceInstanceGUID interface{} `json:"service_instance_guid"`
		Port                interface{} `json:"port"`
		DomainURL           string      `json:"domain_url"`
		SpaceURL            string      `json:"space_url"`
		AppsURL             string      `json:"apps_url"`
		RouteMappingsURL    string      `json:"route_mappings_url"`
	} `json:"entity"`
}

// RoutePage represents the result of cf curl /v2/apps/<guid>/routes output
type RoutePage struct {
	// TotalResults int            `json:"total_results"`
	// TotalPages   int            `json:"total_pages"`
	// PrevURL      string         `json:"prev_url"`
	// NextURL      string         `json:"next_url"`
	Resources []RouteDetails `json:"resources"`
}

// DomainDetails is the Go struct for the curl /v2/shared_domains/<guid> result JSON
type DomainDetails struct {
	Metadata struct {
		GUID      string    `json:"guid"`
		URL       string    `json:"url"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"metadata"`
	Entity struct {
		Name            string      `json:"name"`
		Internal        bool        `json:"internal"`
		RouterGroupGUID interface{} `json:"router_group_guid"`
		RouterGroupType interface{} `json:"router_group_type"`
	} `json:"entity"`
}
