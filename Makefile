BINARY=contacts

# system variables
TARGET_DIR=target
BRANCH ?= $(shell git branch | grep "^\*" | sed 's/^..//')
COMMIT ?= $(shell git rev-parse --short HEAD)
VERSION=$(BRANCH)-$(COMMIT)
BUILDTIME ?= $(shell date -u +%FT%T)

DOCKER_COMPOSE_CMD=docker-compose -p contacts

LINT ?= golangci-lint

# build flags
BUILD_ENV=GOARCH=amd64 CGO_ENABLED=0
LDFLAGS=-ldflags='-w -s -X github.com/davidterranova/contacts/cmd.Version=${VERSION} -X github.com/davidterranova/contacts/cmd.BuildTime=${BUILDTIME}'
BUILD_FLAGS=-a

.PHONY: build
build: clean prepare
	$(BUILD_ENV) $(GOOS) $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(TARGET_DIR)/$(BINARY) .

.PHONY: lint
lint:
	$(LINT) run ./...

.PHONY: lint-fix
lint-fix:
	$(LINT) run --fix ./...

.PHONY: test-unit
test-unit:
	find . -name '.sequence' -type d | xargs rm -rf
	go test ./... -v -count=1 -race -cover

# it requires having a mockgen installed. See: https://github.com/golang/mock
.PHONY: mockgen
mockgen:
	go generate ./...
