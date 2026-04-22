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

# run docker services (db, migrate , app)
docker-up:
	docker compose up --build

docker-up-detached:
	docker compose up -d --build

# stop all docker services
docker-down:
	docker compose down

# run db only
docker-db-up:
	docker compose up -d postgres

# wait
docker-db-wait:
	until docker compose exec -T postgres pg_isready -U "$(POSTGRES_USER)" -d "$(POSTGRES_DB)"; do sleep 1; done

# run db, wait, mirgate
setup: docker-db-up docker-db-wait migrate-up

# clear db
db-clear:
	docker compose down -v

# clear, re-setup db
reset: db-clear setup
