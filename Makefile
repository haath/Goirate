
ARGS = -i -v
OUTPUT = -o build

targets: dep build

dep: Gopkg.toml Gopkg.lock
	dep ensure

build: scanner/main.go
	go build $(ARGS) $(OUTPUT)/scanner scanner/main.go

clean:
	@rm -rf build