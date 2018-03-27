DIRS     ?= $(shell find . -name '*.go' | grep --invert-match 'vendor' | xargs -n 1 dirname | sort --unique)
PKG_NAME ?= aws-creds

BFLAGS ?=
LFLAGS ?=
TFLAGS ?=

VERSION ?=

COVERAGE_PROFILE ?= coverage.out

default: build

build:
	@echo "---> Building"
	go build -o ./bin/$(PKG_NAME) $(BFLAGS)

lint:
	@echo "---> Linting... this might take a minute"
	gometalinter --vendor --tests --deadline=3m $(LFLAGS) $(DIRS)

test:
	@echo "---> Testing"
	go test ./... -coverprofile $(COVERAGE_PROFILE) $(TFLAGS)

enforce:
	@echo "---> Enforcing coverage"
	./scripts/coverage.sh $(COVERAGE_PROFILE)

html:
	@echo "---> Generating HTML coverage report"
	go tool cover -html $(COVERAGE_PROFILE)

release:
	@echo "---> Creating tagged release"
	git tag $(VERSION)
	git push origin
	git push origin --tags

clean:
	@echo "---> Cleaning"
	@rm -rf ./bin

install_tools:
	@echo "--> Installing tools"
	go get -u -v github.com/alecthomas/gometalinter
	gometalinter --install

uninstall_tools:
	@echo "--> Uninstalling tools"
	go clean -i github.com/alecthomas/gometalinter

.PHONY: build lint test enforce html clean install_tools uninstall_tools
