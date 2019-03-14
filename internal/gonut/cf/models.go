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
