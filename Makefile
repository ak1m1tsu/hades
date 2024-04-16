BUILD_COMMIT = $(shell git rev-parse HEAD | cut -c -8)
BUILD_TIME = $(shell TZ=UTC-3 date +%Y%m%d-%H%M)
BUILD_OS = $(shell uname)
BUILD_VERSION ?= $(shell TZ=UTC-3 date -d @$(shell git show -s --format=%ct $(shell git rev-parse HEAD) | head -1) '+%Y%m%d-%H%M')

LD_FLAGS := "-X 'main.BuildCommit=${BUILD_COMMIT}' -X 'main.BuildVersion=${BUILD_VERSION}' -X 'main.BuildTime=${BUILD_TIME}' -X 'main.BuildOS=${BUILD_OS}'"

APPLICATION := hades

.PHONY: build
build:
	@go build -o ./bin/${APPLICATION} -ldflags=${LD_FLAGS} -v ./cmd/${APPLICATION}/main.go

.PHONY: run
run: build
	@./bin/${APPLICATION} --config config/config.yaml

.PHONY: lint
lint:
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.57.2 run --config=.golangci.yaml ./...

.PHONY: init-hook
init-hooks:
	@touch ./.git/hooks/pre-commit
	@echo '#!/bin/sh' > ./.git/hooks/pre-commit
	@echo 'make lint' >> ./.git/hooks/pre-commit
	@chmod +x ./.git/hooks/pre-commit
	@echo 'Git hooks inited'

.PHONY: docker/up
docker/up:
	@docker compose -f docker-compose.yaml up -d --build

.PHONY: docker/down
docker/down:
	@docker compose -f docker-compose.yaml down --rmi local

.PHONY: docker/build
docker/build:
	@docker build
