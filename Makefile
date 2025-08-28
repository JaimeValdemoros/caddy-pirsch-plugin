.PHONY: test deps

test:
	CGO_ENABLED=1 go test -cover -race -p 1 ./pkg/...

deps:
	go get -u -t ./...
	go mod tidy
	go mod vendor
