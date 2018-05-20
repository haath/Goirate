
ARGS = -i -v
COV_MODE = set
TEST_ARGS = -v -covermode=$(COV_MODE)

PKG_LIST := $(shell go list ./... | grep -v /vendor/)

.PHONY: all build clean

all: dep build

test: dep
	-@mkdir -p build
	go fmt $(PKG_LIST)
	go vet -composites=false $(PKG_LIST)

	@echo "mode: $(COV_MODE)" > build/coverage.cov
	@for package in $(PKG_LIST); do \
		go test $(TEST_ARGS) -coverprofile build/tmp.cov $$package ; \
		tail -q -n +2 build/tmp.cov >> build/coverage.cov; \
		rm build/tmp.cov; \
	done
	go tool cover -func=build/coverage.cov

build:
	go build $(ARGS) -o build/gorrent ./cmd

dep: Gopkg.toml Gopkg.lock
	dep ensure

clean:
	-@rm -rf build