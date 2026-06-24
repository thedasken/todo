run:
	@go run .

test:
	@go test -v ./...

build:
	@go build .

clean:
	@rm -f todo

.PHONY: run test build