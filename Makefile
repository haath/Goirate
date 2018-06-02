
OUTPUT := build/goirate

PKG_LIST := $(shell go list ./... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

build: dep compile ## Install dependencies and compile the binary file

install: dep ## Compile and install the binary at $GOPATH/bin
	go install ./cmd/goirate

lint: dep ## Verifies the code through lint, fmt and vet
	@echo "Linting..."
	@golint -set_exit_status $(PKG_LIST)
	@echo "Formatting.."
	@go fmt $(PKG_LIST)
	@echo "Vetting..."
	@go vet -composites=false $(PKG_LIST)

test: dep ## Run unit tests
	@go test -short ${PKG_LIST}

test-cov: dep ## Run unit tests and generate code coverage
	@chmod +x ./scripts/test.sh
	./scripts/test.sh;

compile: ## Compile the binary file
	@go build -i -v -o $(OUTPUT) ./cmd/goirate

dep: Gopkg.toml ## Install dependencies
	@dep ensure
	@go get -u github.com/golang/lint/golint

clean: ## Remove previous build
	@rm -rf build

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "make \033[36m%-30s\033[0m %s\n", $$1, $$2}'