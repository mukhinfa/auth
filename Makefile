include .env

LOCAL_BIN:=$(CURDIR)/bin

LOCAL_MIGRATIONS_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=disable"

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.63.4

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml --fix

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.33.0
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.24.1

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

local-migation-status:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATIONS_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

local-migation-up:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATIONS_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

local-migation-down:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATIONS_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

generate:
	make generate-user_v1-api

generate-user_v1-api:
	mkdir -p pkg/user/v1
	protoc --proto_path api/user/v1 \
	--go_out=pkg/user/v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/user/v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/user/v1/user.proto