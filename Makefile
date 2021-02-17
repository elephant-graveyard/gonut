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

.PHONY: all clean generated-code lint vet misspell shellcheck unit-test test build

all: test build

clean:
	@rm -rf dist
	@assets/sample-apps-src/java/clean.sh
	@assets/sample-apps-src/binary/clean.sh
	@pina-golada cleanup
	@go clean -i -cache -testcache $(shell go list ./...)

fmts:
	@find . -type f -name '*.go' -print0 | xargs -0 gofmt -s -w -e -l
	@shfmt -s -w -i 2 -ci scripts/*.sh

assets/sample-apps/binary/binary:
	@echo "Compiling example binary app ..."
	@assets/sample-apps-src/binary/compile.sh
	@echo

assets/sample-apps/java/App.jar:
	@echo "Compiling example Java app ..."
	@assets/sample-apps-src/java/compile.sh
	@echo

generated-code: assets/sample-apps/binary/binary assets/sample-apps/java/App.jar
	@pina-golada generate

lint:
	@golint --set_exit_status ./...

vet:
	@go vet ./...

misspell:
	@find . -type f \( -name "*.go" -o -name "*.md" \) -print0 | xargs -0 misspell -error

shellcheck:
	@shellcheck scripts/*.sh

unit-test:
	@ginkgo -r \
	--randomizeAllSpecs \
	--randomizeSuites \
	--failOnPending \
	--nodes=4 \
	--compilers=2 \
	--race \
	--trace

test: generated-code lint vet misspell shellcheck unit-test
