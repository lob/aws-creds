BIN_DIR ?= ./bin
PKG_NAME ?= aws-creds
GO_TOOLS := \
	github.com/git-chglog/git-chglog/cmd/git-chglog \

VERSION ?=

COVERAGE_PROFILE ?= coverage.out

default: build

.PHONY: build
build:
	@echo "---> Building"
	go build -ldflags "-w -s" -o ./bin/$(PKG_NAME) ./cmd/aws-creds

.PHONY: clean
clean:
	@echo "---> Cleaning"
	go clean
	rm -rf ./bin

.PHONY: enforce
enforce:
	@echo "---> Enforcing coverage"
	./scripts/coverage.sh $(COVERAGE_PROFILE)

.PHONY: html
html:
	@echo "---> Generating HTML coverage report"
	go tool cover -html $(COVERAGE_PROFILE)

.PHONY: install
install:
	@echo "---> Installing dependencies"
	go mod download

.PHONY: lint
lint:
	@echo "---> Linting"
	$(BIN_DIR)/golangci-lint run

.PHONY: release
release:
	@echo "---> Creating tagged release"
	git tag $(VERSION)
	git push origin
	git push origin --tags

.PHONY: setup
setup: install
	@echo "--> Setting up tools"
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(BIN_DIR) v1.20.0
	go get $(GO_TOOLS) && GOBIN=$$(realpath $(BIN_DIR)) go install $(GO_TOOLS)

.PHONY: test
test:
	@echo "---> Testing"
	go test ./... -race -coverprofile $(COVERAGE_PROFILE)
