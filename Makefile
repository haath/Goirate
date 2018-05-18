
ARGS = -i -v
OUTPUT = -o build

GO_FILES := $(shell go list ./... | grep -v /vendor/)

targets: build

test: dep
	-@mkdir -p build
	go fmt $(GO_FILES)
	go vet -composites=false $(GO_FILES)

	go test -v -coverprofile build/scanner.testCoverage.txt ./scanner
	go test -v -coverprofile build/shared.testCoverage.txt ./shared

build: dep
	go build $(ARGS) $(OUTPUT)/scanner ./scanner

dep: Gopkg.toml Gopkg.lock
	dep ensure

clean:
	-@rm -rf build