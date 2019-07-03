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

package assets

import (
	"github.com/homeport/pina-golada/pkg/files"
)

// Provider is the provider instance to access pina-goladas asset framework
var Provider ProviderInterface

// ProviderInterface is the interface that provides the asset file
// @pgl(package=assets&injector=Provider)
type ProviderInterface interface {

	// GoSampleApp returns the directory containing the Golang sample app
	// @pgl(asset=/assets/sample-apps/golang/&compressor=tar)
	GoSampleApp() (directory files.Directory, e error)

	// PythonSampleApp returns the directory containing the Python sample app
	// @pgl(asset=/assets/sample-apps/python/&compressor=tar)
	PythonSampleApp() (directory files.Directory, e error)

	// PHPSampleApp returns the directory containing the PHP sample app
	// @pgl(asset=/assets/sample-apps/php/&compressor=tar)
	PHPSampleApp() (directory files.Directory, e error)

	// StaticfileSampleApp returns the directory containing the Staticfile sample app
	// @pgl(asset=/assets/sample-apps/staticfile/&compressor=tar)
	StaticfileSampleApp() (directory files.Directory, e error)

	// SwiftSampleApp returns the directory containing the Swift sample app
	// @pgl(asset=/assets/sample-apps/swift/&compressor=tar)
	SwiftSampleApp() (directory files.Directory, e error)

	// NodeJSSampleApp returns the directory containing the NodeJS sample app
	// @pgl(asset=/assets/sample-apps/nodejs/&compressor=tar)
	NodeJSSampleApp() (directory files.Directory, e error)

	// RubySampleApp returns the directory containing the Ruby sample app
	// @pgl(asset=/assets/sample-apps/ruby/&compressor=tar)
	RubySampleApp() (directory files.Directory, e error)

	// BinarySampleApp returns the directory containing the Binary sample app
	// @pgl(asset=/assets/sample-apps/binary/&compressor=tar)
	BinarySampleApp() (directory files.Directory, e error)
}
