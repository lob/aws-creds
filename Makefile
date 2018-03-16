DIRS     ?= $(shell find . -name '*.go' | grep --invert-match 'vendor' | xargs -n 1 dirname | sort --unique)
PKG_NAME ?= aws-creds

BFLAGS ?=
LFLAGS ?=
TFLAGS ?=

COVERAGE_PROFILE ?= coverage.out

default: build

build:
	@echo "---> Building"
	go build -o ./bin/$(PKG_NAME) $(BFLAGS)

lint:
	@echo "---> Linting... this might take a minute"
	gometalinter --vendor --tests --deadline=2m $(LFLAGS) $(DIRS)

test:
	@echo "---> Testing"
	go test ./... -coverprofile $(COVERAGE_PROFILE) $(TFLAGS)

enforce:
	@echo "---> Enforcing coverage"
	./scripts/coverage.sh $(COVERAGE_PROFILE)

html:
	@echo "---> Generating HTML coverage report"
	go tool cover -html $(COVERAGE_PROFILE)

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
