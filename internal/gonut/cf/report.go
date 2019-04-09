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
	"strings"
	"time"
)

// PushReport encapsules details of a Cloud Foundry push command
type PushReport struct {
	InitStart      time.Time
	CreatingStart  time.Time
	UploadingStart time.Time
	StagingStart   time.Time
	StartingStart  time.Time
	PushEnd        time.Time

	buildpack *BuildpackDetails
	stack     *StackDetails
}

// InitTime is the time it takes to initialise the Cloud Foundry app push setup
func (report PushReport) InitTime() time.Duration {
	return report.CreatingStart.Sub(report.InitStart)
}

// CreatingTime is the time it takes to create the app in Cloud Foundry
func (report PushReport) CreatingTime() time.Duration {
	return report.UploadingStart.Sub(report.CreatingStart)
}

// UploadingTime is the time it takes to upload the app bits to Cloud Foundry
func (report PushReport) UploadingTime() time.Duration {
	return report.StagingStart.Sub(report.UploadingStart)
}

// StagingTime is the time it takes to stage (compile) the app in Cloud Foundry
func (report PushReport) StagingTime() time.Duration {
	return report.StartingStart.Sub(report.StagingStart)
}

// StartingTime is the time it takes to start the compiled app in Cloud Foundry
func (report PushReport) StartingTime() time.Duration {
	return report.PushEnd.Sub(report.StartingStart)
}

// ElapsedTime is the overall elapsed time it takes to push an app in Cloud Foundry
func (report PushReport) ElapsedTime() time.Duration {
	return report.PushEnd.Sub(report.InitStart)
}

// Buildpack provides the name of the buildpack used (if detectable)
func (report PushReport) Buildpack() string {
	if report.buildpack != nil {
		return report.buildpack.Entity.Name
	}

	return "(unknown)"
}

// Stack provides the name of the stack used (if detectable)
func (report PushReport) Stack() string {
	if report.stack != nil {
		return fmt.Sprintf("%s (%s)",
			report.stack.Entity.Description,
			report.stack.Entity.Name,
		)
	}

	return "(unknown)"
}

// ParseUpdate parses a line from the CF CLI push output
func (report *PushReport) ParseUpdate(text string) string {
	switch {
	case strings.HasPrefix(text, "Creating app"):
		report.CreatingStart = time.Now()
		return "Creating"

	case strings.HasPrefix(text, "Uploading files") || strings.HasPrefix(text, "Uploading app files from"):
		report.UploadingStart = time.Now()
		return "Uploading"

	case strings.HasPrefix(text, "Staging app") || strings.HasPrefix(text, "Staging..."):
		report.StagingStart = time.Now()
		return "Staging"

	case strings.HasPrefix(text, "Waiting for app to start...") || strings.HasPrefix(text, "Successfully destroyed container"):
		report.StartingStart = time.Now()
		return "Starting"

	case strings.HasPrefix(text, "Deleting app"):
		return "Deleting"
	}

	return ""
}
