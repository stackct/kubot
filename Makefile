PROJECT_BUILD_VERSION:=$(shell cat gen/VERSION | tr -d '\n')$(shell git rev-parse --short HEAD)
PROJECT_BUILD_DATE="$(shell date -u +%FT%T.000Z)"
LDFLAGS=-ldflags=all="-X 'main.version=$(PROJECT_BUILD_VERSION)' -X 'main.buildDate=$(PROJECT_BUILD_DATE)'"

default: dist

dist: test dist-linux

dist-linux:
	@go get -u -d ./... && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kubot $(LDFLAGS) .

dep:
	@go get -v -u -d ./...

dep-test:
	@go get ./...

update:
	@go get -u ./...

install: dep
	@go build -a -o $(GOHOME)/bin/$(BINARY) $(LDFLAGS)

clean:
	@rm $(BINARY)

test: dep-test
	@go test -v -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: all