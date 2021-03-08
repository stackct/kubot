default: test

update:
	@go get -u ./...

test: update
	@go test -v -coverprofile=coverage.txt -covermode=atomic ./...
