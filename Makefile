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

GRPC_DST_DIR=.
GRPC_SRC_DIR=./internal/adapters/grpc

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
test-unit: compose-up
	find . -name '.sequence' -type d | xargs rm -rf
	go test ./... -v -count=1 -race -cover

# it requires having a mockgen installed. See: https://github.com/golang/mock
.PHONY: mockgen
mockgen:
	go generate ./...

.PHONY: gqlgen
gqlgen:
	go get github.com/99designs/gqlgen@v0.17.35
	go run github.com/99designs/gqlgen generate

.PHONY: grpcgen
grpcgen:
	protoc -I=$(GRPC_SRC_DIR) --go_out=$(GRPC_DST_DIR) $(GRPC_SRC_DIR)/contacts.proto
	protoc --go_out=$(GRPC_DST_DIR) --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative $(GRPC_SRC_DIR)/contacts.proto

.PHONY: compose-up
compose-up:
	docker-compose up -d

.PHONY: compose-down
compose-down:
	docker-compose down

.PHONY: gen-migration
gen-migration:
	migrate create -ext sql -dir pkg/pg/migrations -seq $(name)

.PHONY: migrate-up
migrate-up:
	migrate -path pkg/pg/migrations -database "$(MIGRATE_DB_CONN_STRING)&sslmode=disable" -verbose up

.PHONY: migrate-down
migrate-down:
	migrate -path pkg/pg/migrations -database "$(MIGRATE_DB_CONN_STRING)&sslmode=disable" -verbose down
