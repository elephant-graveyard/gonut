# Copyright © 2019 The Homeport Team
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
# THE SOFTWARE.

.PHONY: all clean prereqs lint vet misspell shellcheck unit-test test build

all: test build

clean:
	@rm -rf binaries
	@assets/sample-apps-src/java/clean.sh
	@assets/sample-apps-src/binary/clean.sh
	@pina-golada cleanup
	@GO111MODULE=on go clean -i -cache -testcache $(shell go list ./...)

prereqs:
	@assets/sample-apps-src/java/compile.sh
	@assets/sample-apps-src/binary/compile.sh

lint:
	@scripts/go-lint.sh

vet:
	@scripts/go-vet.sh

misspell:
	@scripts/misspell.sh

shellcheck:
	@scripts/shellcheck.sh

unit-test:
	@scripts/test.sh

test: prereqs lint vet misspell shellcheck unit-test

build: prereqs
	@scripts/build.sh
