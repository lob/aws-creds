DIRS     ?= $(shell find . -name '*.go' | grep --invert-match 'vendor' | xargs -n 1 dirname | sort --unique)
PKG_NAME ?= aws-creds

BFLAGS ?=
LFLAGS ?=
TFLAGS ?=

TS_METHOD ?=
TS_SUITE  ?=

default: build

build:
	@echo "---> Building"
	go build -o ./bin/$(PKG_NAME) $(BFLAGS)

build_all:
	@echo "---> Building: (darwin, amd64)"
	GOOS=darwin GOARCH=amd64 go build -v -o ./bin/$(PKG_NAME)_osx $(BFLAGS)
	@echo "---> Building: (windows, amd64)"
	GOOS=windows GOARCH=amd64 go build -v -o ./bin/$(PKG_NAME)_win.exe $(BFLAGS)

lint:
	@echo "---> Linting... this might take a minute"
	gometalinter --vendor --tests --deadline=2m $(LFLAGS) $(DIRS)

test:
	@echo "---> Testing"
	go test ./... -cover $(TFLAGS) | tee ./coverage.out

test_only:
	@echo "---> Running unit tests only, no coverage"
	go test ./... -short $(TFLAGS)

test_unit:
	@echo "---> Testing unit tests"
	go test ./... -short -cover $(TFLAGS)

test_integration:
	@echo "---> Testing integration tests"
	go test ./... -run Integration $(TFLAGS)

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

.PHONY: build build_all lint test test_only clean install_tools uninstall_tools
