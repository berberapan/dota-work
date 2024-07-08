.PHONY: build run test clean

build:
	@go build -o bin/app cmd/main.go

run: build
	@./bin/app

test:
	@go test -v ./test/... 

clean:
	@rm -f bin/app
