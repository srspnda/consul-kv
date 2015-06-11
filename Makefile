NAME = $(shell $(basename "$(pwd)"))
VERSION = $(shell grep -oE "[0-9]\.[0-9]\.[0-9]" main.go)
DEPS = $(shell go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)

GOXOS = "linux darwin windows"
GOXOUT = "build/{{.Dir}}_$(VERSION)_{{.OS}}_{{.Arch}}/$(NAME)"

all: deps build

deps:
	go get -d -v ./...
	echo $(DEPS) | xargs -n1 go get -d

build: deps
	@mkdir -p bin/
	go build -o bin/$(NAME)

test:
	go test ./...
	go vet ./...

gox: build test
	rm -rf build
	@mkdir -p build
	gox -os=$(GOXOS) -output=$(GOXOUT)

.PHONY: all deps build test gox
