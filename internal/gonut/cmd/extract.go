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
	"os"
	"path/filepath"

	"github.com/homeport/gonut/internal/gonut/nok"
	"github.com/homeport/pina-golada/pkg/files"
	"github.com/spf13/cobra"
)

// extractCmd represents the push command
var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extracts the sample apps to disk",
	Long:  `Extracts the sample apps to the local disk.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		baseDir := "."

		for _, sampleApp := range sampleApps {
			path := filepath.Join(baseDir, "SampleApps", sampleApp.caption)
			if err := os.MkdirAll(path, os.FileMode(0755)); err != nil {
				return err
			}

			directory, err := sampleApp.assetFunc()
			if err != nil {
				return &nok.ErrorWithDetails{
					Caption: "failed to load sample app assets into memory",
					Details: err.Error(),
				}
			}

			if err := files.WriteToDisk(directory, path, true); err != nil {
				return &nok.ErrorWithDetails{
					Caption: "failed to write sample app files to disk",
					Details: err.Error(),
				}
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(extractCmd)
}
