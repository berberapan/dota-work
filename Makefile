.PHONY: build run test clean front-install front-build front-run run-all

# Backend

build:
	@go build -o bin/app cmd/main.go

run: build
	@./bin/app

test:
	@go test -v ./test/... 

clean:
	@rm -f bin/app
	@rm -rf web/dist

# Frontend

front-install:
	@cd web && npm install

front-build:
	@cd web && npm run build

front-run:
	@cd web && npm run dev

# Combined

run-all:
	@(make run & make front-run)
