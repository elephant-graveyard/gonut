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

package cf_test

import (
	"bufio"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/homeport/gonut/internal/gonut/cf"
)

func linefeeder(path string, fn func(text string)) {
	file, err := os.Open(path)
	Expect(err).ToNot(HaveOccurred())
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fn(scanner.Text())
	}

	Expect(scanner.Err()).ToNot(HaveOccurred())
}

func createMockReport(path string) *PushReport {
	report := &PushReport{
		AppName:   "the-app-name",
		InitStart: time.Now(),
	}

	linefeeder(path, func(text string) {
		report.ParseUpdate(text)
	})

	report.PushEnd = time.Now()

	return report
}

var _ = Describe("Cloud Foundry push report", func() {
	Context("Parse Cloud Foundry push output", func() {
		var (
			unset = time.Time{}
		)

		It("should parse cloud controller api version 2.92.0 style logs", func() {
			report := createMockReport("../../../assets/test/cf-push/api-2.92.0/push-and-delete.log")

			Expect(report.InitStart).ToNot(BeEquivalentTo(unset))
			Expect(report.CreatingStart).ToNot(BeEquivalentTo(unset))
			Expect(report.UploadingStart).ToNot(BeEquivalentTo(unset))
			Expect(report.StagingStart).ToNot(BeEquivalentTo(unset))
			Expect(report.StartingStart).ToNot(BeEquivalentTo(unset))
		})

		It("should parse cloud controller api version 2.106.0 style logs", func() {
			report := createMockReport("../../../assets/test/cf-push/api-2.106.0/push-and-delete.log")

			Expect(report.InitStart).ToNot(BeEquivalentTo(unset))
			Expect(report.CreatingStart).ToNot(BeEquivalentTo(unset))
			Expect(report.UploadingStart).ToNot(BeEquivalentTo(unset))
			Expect(report.StagingStart).ToNot(BeEquivalentTo(unset))
			Expect(report.StartingStart).ToNot(BeEquivalentTo(unset))
		})

		It("should parse cloud controller api version 2.128.0 style logs", func() {
			report := createMockReport("../../../assets/test/cf-push/api-2.128.0/push-and-delete.log")

			Expect(report.InitStart).ToNot(BeEquivalentTo(unset))
			Expect(report.CreatingStart).ToNot(BeEquivalentTo(unset))
			Expect(report.UploadingStart).ToNot(BeEquivalentTo(unset))
			Expect(report.StagingStart).ToNot(BeEquivalentTo(unset))
			Expect(report.StartingStart).ToNot(BeEquivalentTo(unset))
		})

		It("should parse cloud controller api version 2.133.0 style logs", func() {
			report := createMockReport("../../../assets/test/cf-push/api-2.133.0/push-and-delete.log")

			Expect(report.InitStart).ToNot(BeEquivalentTo(unset))
			Expect(report.CreatingStart).ToNot(BeEquivalentTo(unset))
			Expect(report.UploadingStart).ToNot(BeEquivalentTo(unset))
			Expect(report.StagingStart).ToNot(BeEquivalentTo(unset))
			Expect(report.StartingStart).ToNot(BeEquivalentTo(unset))
		})

		It("should parse unknown cloud controller api style logs without issues", func() {
			report := createMockReport("../../../assets/test/cf-push/api-unknown/push-and-delete.log")

			Expect(report.InitStart).ToNot(BeEquivalentTo(unset))
			Expect(report.PushEnd).ToNot(BeEquivalentTo(unset))
		})
	})
})
