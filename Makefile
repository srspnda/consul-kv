VERSION=0.1.0

all: build

build:
	@mkdir -p bin/
	go build -o bin/consul-kv
	cp bin/consul-kv ${GOPATH}/bin

release:
	@mkdir -p bin/
	gox -os="linux windows darwin" -output="bin/{{.Dir}}_${VERSION}_{{.OS}}_{{.Arch}}"

.PHONY: all build release
