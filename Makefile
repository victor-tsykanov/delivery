include .env

LOCAL_BIN:=$(CURDIR)/bin
DB_DSN="host=localhost port=$(DB_PORT) dbname=$(DB_DATABASE) user=$(DB_USER) password=$(DB_PASSWORD) sslmode=$(DB_SSL_MODE)"

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.5
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1

generate:
	rm -rf pkg/servers
	rm -rf mocks
	mkdir -p pkg/servers
	go tool mockery
	go tool oapi-codegen -config configs/server.cfg.yaml https://gitlab.com/microarch-ru/ddd-in-practice/system-design/-/raw/main/services/delivery/contracts/openapi.yml

test:
	go test ./...

build:
	go build -o bin/server cmd/server.go

lint:
	go tool golangci-lint run ./...

fix:
	go tool golangci-lint run --fix ./...

migrate-up:
	go tool goose -dir ./migrations postgres ${DB_DSN} up -v

migrate-down:
	go tool goose -dir ./migrations postgres ${DB_DSN} down -v

start-dev-server:
	go tool air --build.cmd "go build -o bin/main cmd/main.go" --build.bin "./bin/main"

test-continuously:
	go tool air --build.cmd "go test ./..." --build.bin ""
