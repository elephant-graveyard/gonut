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
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/gonvenience/bunt"
	"github.com/gonvenience/neat"
	"github.com/gonvenience/wrap"
	"github.com/homeport/gonut/internal/gonut/nok"
)

// rootCmd represents the base command when called without any sub-commands
var rootCmd = &cobra.Command{
	Use:   "gonut",
	Short: "Portable tool to help you verify whether you can push a sample app to a Cloud Foundry.",
	Long: bunt.Sprintf(`
Gonut is a portable tool to help you verify whether you can push a sample app to
a Cloud Foundry instance. It will push a sample app to Cloud Foundry and delete
it afterwards. The apps are embedded into the gonut binary.

It is written in Golang and uses CornflowerBlue{~https://github.com/homeport/pina-golada~} to
include arbitrary sample app data in the application binary.`),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		ExitGonut(err)
	}
}

// ExitGonut leaves gonut in case of an unresolvable error situation
func ExitGonut(reason interface{}) {
	var (
		headline string
		content  string
	)

	switch typed := reason.(type) {
	case *nok.ErrorWithDetails:
		headline = bunt.Sprintf("*Error:* _%s_", typed.Caption)
		content = fmt.Sprintf("%s\n\n", typed.Details)

	case wrap.ContextError:
		headline = bunt.Sprintf("*Error:* _%s_", typed.Context())
		content = typed.Cause().Error()

	default:
		headline = "Error occurred"
		content = fmt.Sprint(reason)
	}

	neat.Box(os.Stderr,
		headline, strings.NewReader(content),
		neat.HeadlineColor(bunt.Coral),
		neat.ContentColor(bunt.DimGray),
	)

	os.Exit(1)
}
