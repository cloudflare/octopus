#!/usr/bin/make

build:
	cd cmd/octopus && go build

test:
	@echo "Running tests..."
	@mkdir -p .cover
	@go test -race -vet=all -cover -coverprofile=.cover/coverage.txt -coverpkg=./pkg/... ./...
	@go tool cover -html ".cover/coverage.txt" -o .cover/all.html

lint:
	@echo "Linting code..."
	@golangci-lint run

proto-go:
	protoc \
		-Ivendor/github.com/bio-routing/bio-rd/ \
		--proto_path=proto \
		--go_out=proto/octopus/ --go_opt=paths=source_relative \
		--go-grpc_out=proto/octopus/ --go-grpc_opt=paths=source_relative,require_unimplemented_servers=false \
		octopus.proto

all: proto-go build lint test

run-octopus: all
	cmd/octopus/octopus -mock-connectors
