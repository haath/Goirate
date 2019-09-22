
OUTPUT := build/goirate

PKG_LIST := $(shell go list ./... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

BUILD_FLAGS		= -i -v -o $(OUTPUT)
GCC_FLAGS		= --ldflags "-linkmode external -extldflags -static -s -w"
GCC_FLAGS_WIN	= --ldflags "-extldflags -static"

GOX_FLAGS		= -ldflags "-X main.Version=$(CI_COMMIT_TAG) -s -w" -tags netgo
GOX_ARCHS		= -osarch="darwin/amd64" -os="linux" -os="windows" -os="solaris" 
GOX_OUTPUT		= "build/goirate.{{.OS}}.{{.Arch}}"

build: dep patch ## Install dependencies and statitcally compile the binary file
	packr build $(GCC_FLAGS) $(BUILD_FLAGS) ./cmd/goirate
	@packr clean

build-win64: dep patch ## Installs dependencies and statitcally compiles the binary file on 64-bit windows
	packr build $(GCC_FLAGS_WIN) $(BUILD_FLAGS) ./cmd/goirate
	@packr clean

cross-compile: dep patch ## Installs dependencies and statitcally cross-compiles the binary using gox
	go get github.com/mitchellh/gox
	packr
	export CGO_ENABLED=0 ;\
	gox $(GOX_FLAGS) $(GOX_ARCHS) -output $(GOX_OUTPUT) ./cmd/goirate
	@packr clean

install: dep patch ## Compiles and install the binary at $GOPATH/bin
	packr install ./cmd/goirate
	@packr clean

lint: dep ## Verifies the code through lint, fmt and vet
	@echo "Linting..."
	@golint -set_exit_status $(PKG_LIST)
	@echo "Formatting.."
	@go fmt $(PKG_LIST)
	@echo "Vetting..."
	@go vet -composites=false $(PKG_LIST)

fmt: ## Runs go fmt on each of the packages
	gofmt -s -w ./cmd
	gofmt -s -w ./pkg

test: dep ## Runs unit tests
	@go test -short ${PKG_LIST}

test-cov: dep ## Rusn unit tests and generate code coverage
	@export GOIRATE_DEBUG=false
	@./scripts/test.sh;

compile: ## Compiles the binary file
	@packr build -i -v -o $(OUTPUT) ./cmd/goirate

dep: Gopkg.toml ## Installs package dependencies with dep ensure
	@dep ensure

install-tools: ## Installs dep, packr and golint
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/gobuffalo/packr/...
	go get -u golang.org/x/lint/golint

clean: ## Remove previous build
	@rm -rf build

patch:
	@./scripts/patch.sh;

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "make \033[36m%-30s\033[0m %s\n", $$1, $$2}'
	