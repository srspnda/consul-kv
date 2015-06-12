NAME = $(shell awk -F\" '/^const Name/ { print $$2 }' main.go)
VERSION = $(shell grep -oE "[0-9]\.[0-9]\.[0-9]" main.go)
DEPS = $(shell go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)

GOXOS = "linux darwin windows"
GOXOUT = "build/{{.Dir}}_$(VERSION)_{{.OS}}_{{.Arch}}/$(NAME)"

all: deps build

deps:
	go get -d -v ./...
	echo $(DEPS) | xargs -n1 go get -d

build:
	@mkdir -p bin/
	go build -o bin/$(NAME)
	cp bin/$(NAME) $(GOPATH)/bin

test: deps
	go test ./...
	go vet ./...

xcompile: deps test
	rm -rf build/
	@mkdir -p build
	gox -os=$(GOXOS) -output=$(GOXOUT)

package: xcompile
	$(eval FILES := $(shell ls build))
	@mkdir -p build/tgz
	for f in $(FILES); do \
		(cd $(shell pwd)/build && tar -czvf tgz/$$f.tar.gz $$f); \
		echo $$f; \
	done

.PHONY: all deps build test xcompile package
