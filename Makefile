BIN_DIR ?= ./bin
PKG_NAME ?= aws-creds
GO_TOOLS := \
	github.com/git-chglog/git-chglog/cmd/git-chglog \

COVERAGE_PROFILE ?= coverage.out

# Override version by setting the VERSION environment variable
VERSION ?=
ifneq ($(strip $(VERSION)),)
LDFLAGS ?= "-w -s -X github.com/lob/aws-creds/pkg/cmd.version=$(VERSION)"
else
LDFLAGS ?= "-w -s"
endif

default: build

.PHONY: build
build:
	@echo "---> Building"
	go build -ldflags $(LDFLAGS) -o $(BIN_DIR)/$(PKG_NAME) ./cmd/aws-creds

.PHONY: clean
clean:
	@echo "---> Cleaning"
	go clean
	rm -rf $(BIN_DIR)

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
	go get $(GO_TOOLS) && GOBIN=$$(cd $(BIN_DIR) && pwd) go install $(GO_TOOLS)

.PHONY: test
test:
	@echo "---> Testing"
	go test ./... -race -coverprofile $(COVERAGE_PROFILE)
