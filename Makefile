all: build

build:
	@mkdir -p bin/
	go build -o bin/consul-kv
	cp bin/consul-kv ${GOPATH}/bin

.PHONY: all build
