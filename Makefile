include .env
export

# make: make migrate-down
migrate-up:
	migrate -path internal/migrations -database "$(DATABASE_URL)" up

# make migrate-up
migrate-down:
	migrate -path internal/migrations -database "$(DATABASE_URL)" down 1

# make mock-service
mock-service:
	mockery --name Service --dir internal/service/bank --output internal/mocks --outpkg mocks

# make mock-repo
mock-repo:
	mockery --name Repository --dir internal/repository/bank --output internal/mocks --outpkg mocks

# make mock-gateway
mock-gateway:
	mockery --name Gateway --dir internal/blockchain --output internal/mocks --outpkg mocks

# make mock-clean-svc
mock-clean-svc:
	rm -f internal/mocks/Service.go

# make mock-clean-repo
mock-clean-repo:
	rm -f internal/mocks/Repository.go

# make mock-clean-gateway
mock-clean-gateway:
	rm -f internal/mocks/Gateway.go

mock-all: mock-service mock-repo mock-gateway

mock-refresh: mock-clean-svc mock-clean-repo mock-clean-gateway mock-all

test-handler:
	GOCACHE=/tmp/hela-bank-sc-go-cache go test ./internal/handler/bank

test-service:
	GOCACHE=/tmp/hela-bank-sc-go-cache go test ./internal/service/bank

test-repo:
	GOCACHE=/tmp/hela-bank-sc-go-cache go test ./internal/repository/bank

test-router:
	GOCACHE=/tmp/hela-bank-sc-go-cache go test ./internal/router

test-bank: test-handler test-service test-repo test-router

test-all:
	GOCACHE=/tmp/hela-bank-sc-go-cache go test ./...
