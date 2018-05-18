
ARGS = -i -v
TEST_ARGS = -v -covermode=set

PKG_LIST := $(shell go list ./... | grep -v /vendor/)

targets: build

test: dep
	-@mkdir -p build
	go fmt $(PKG_LIST)
	go vet -composites=false $(PKG_LIST)

	@echo "mode: set" > build/testCoverage.cov
	for package in $(PKG_LIST); do \
		go test $(TEST_ARGS) -coverprofile build/tmp.cov $$package ; \
		tail -q -n +2 build/tmp.cov >> build/testCoverage.cov; \
		rm build/tmp.cov; \
	done
	go tool cover -func=build/testCoverage.cov

build: dep
	go build $(ARGS) -o build/scanner ./scanner

dep: Gopkg.toml Gopkg.lock
	dep ensure

clean:
	-@rm -rf build