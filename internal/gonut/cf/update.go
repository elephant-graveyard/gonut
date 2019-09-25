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
	"fmt"
	"image/color"

	"github.com/gonvenience/bunt"
	"github.com/lucasb-eyer/go-colorful"

	"github.com/gonvenience/neat"
	"github.com/tcnksm/go-latest"
)

var version string
var goos string

// NotifyLatestRelease compares the local version of gonut with the latest release
// on GitHub and prints a notification in case the local version is outdated.
func NotifyLatestRelease() error {
	githubTag := &latest.GithubTag{
		Owner:      "homeport",
		Repository: "gonut",
	}

	if res, err := latest.Check(githubTag, version); err == nil && res.Outdated {
		updateCommand := ""
		switch goos {
		case "darwin":
			updateCommand = "brew install homeport/tap/gonut"
		default:
			updateCommand = "curl -fsL http://ibm.biz/Bd2t2v | bash"
		}

		headline := "New version of gonut available!"
		table, _ := neat.Table(releaseNotificationTable(version, res.Current))
		content := bunt.Sprintf("%s\nRun `YellowGreen{%s}` to update!",
			table,
			updateCommand,
		)

		gonutBlue, _ := colorful.MakeColor(color.RGBA{77, 173, 233, 255})
		formatted := neat.ContentBox(headline, content, neat.HeadlineColor(gonutBlue))

		fmt.Printf("\n%s", formatted)
	}
	return nil
}

func releaseNotificationTable(installed string, current string) [][]string {
	rows := [][]string{
		[]string{
			bunt.Sprintf("*Release:*"),
			bunt.Sprintf("Crimson{%s} → YellowGreen{%s}", installed, current),
		},
		[]string{
			bunt.Sprintf("*Changelog:*"),
			bunt.Sprintf("CornflowerBlue{~https://github.com/homeport/gonut/releases/tag/v%s~}", current),
		},
	}

	return rows
}
