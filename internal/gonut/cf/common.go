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

// HumanReadableDuration returns a human readable version of the time duration
func HumanReadableDuration(duration time.Duration) string {
	if duration < time.Second {
		return "less than a second"
	}

	seconds := int(duration.Seconds())
	minutes := 0
	hours := 0

	if seconds >= 60 {
		minutes = seconds / 60
		seconds = seconds % 60

		if minutes >= 60 {
			hours = minutes / 60
			minutes = minutes % 60
		}
	}

	parts := []string{}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%d h", hours))
	}

	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%d min", minutes))
	}

	if seconds > 0 {
		parts = append(parts, fmt.Sprintf("%d sec", seconds))
	}

	return strings.Join(parts, " ")
}
