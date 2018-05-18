
ARGS = -i -v
OUTPUT = -o build

GO_FILES := $(shell go list ./... | grep -v /vendor/)

targets: build

test: dep
	@mkdir build
	go fmt $(GO_FILES)
	go vet -composites=false $(GO_FILES)
	go test $(GO_FILES) -v -coverprofile build/.testCoverage.txt

build: dep
	go build $(ARGS) $(OUTPUT)/scanner scanner/main.go

dep: Gopkg.toml Gopkg.lock
	dep ensure

clean:
	@rm -rf build