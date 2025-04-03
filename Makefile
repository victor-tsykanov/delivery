include .env

LOCAL_BIN:=$(CURDIR)/bin
DB_DSN="host=localhost port=$(DB_PORT) dbname=$(DB_DATABASE) user=$(DB_USER) password=$(DB_PASSWORD) sslmode=$(DB_SSL_MODE)"

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.6
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1

generate: generate-mocks generate-delivery-api generate-geo-api generate-queues-events

generate-mocks:
	rm -rf mocks
	go tool mockery

generate-delivery-api:
	rm -rf pkg/servers
	mkdir -p pkg/servers
	go tool oapi-codegen -config configs/server.cfg.yaml https://gitlab.com/microarch-ru/ddd-in-practice/system-design/-/raw/main/services/delivery/contracts/openapi.yml

generate-geo-api:
	wget https://gitlab.com/microarch-ru/ddd-in-practice/system-design/-/raw/main/services/geo/contracts/contract.proto -O api/geo/contract.proto
	rm -rf pkg/clients/geopb
	mkdir -p pkg/clients/geopb
	protoc --proto_path api/geo \
		--go_out=pkg/clients/geopb --go_opt=paths=source_relative \
		--plugin=protoc-gen-go=bin/protoc-gen-go \
		--go-grpc_out=pkg/clients/geopb --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
		api/geo/contract.proto

generate-queues-events:
	wget https://gitlab.com/microarch-ru/ddd-in-practice/system-design/-/raw/main/services/basket/contracts/basket_confirmed.proto -O api/basket/queues/basket_confirmed.proto
	protoc --go_out=./pkg --plugin=protoc-gen-go=bin/protoc-gen-go api/basket/queues/basket_confirmed.proto
	wget https://gitlab.com/microarch-ru/ddd-in-practice/system-design/-/raw/main/services/delivery/contracts/order_status_changed.proto -O api/delivery/queues/order_status_changed.proto
	protoc --go_out=./pkg --plugin=protoc-gen-go=bin/protoc-gen-go api/delivery/queues/order_status_changed.proto

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
