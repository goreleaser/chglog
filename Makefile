SOURCE_FILES?=./...
TEST_PATTERN?=.
TEST_OPTIONS?=
TEST_TIMEOUT?=5m

# Install all the build and lint dependencies
setup:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh
	go mod tidy
	git config core.hooksPath .githooks
.PHONY: setup

test:
	go test $(TEST_OPTIONS) -v -failfast -race -coverpkg=./... -covermode=atomic -coverprofile=coverage.out $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=$(TEST_TIMEOUT)
.PHONY: test

cover: test
	go tool cover -html=coverage.out
.PHONY: cover

fmt:
	go mod tidy
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done
.PHONY: fmt

lint: check
	./bin/golangci-lint run --exclude-use-default=false --fix -v
.PHONY: check

ci: build lint test
.PHONY: ci

build:
	go build -o nfpm ./cmd/chglog/main.go
.PHONY: build

deps:
	go get -u
	go mod tidy
	go mod verify
.PHONY: deps

todo:
	@grep \
		--exclude-dir=vendor \
		--exclude-dir=node_modules \
		--exclude-dir=bin \
		--exclude=Makefile \
		--text \
		--color \
		-nRo -E ' TODO:.*|SkipNow' .
.PHONY: todo

.DEFAULT_GOAL := build