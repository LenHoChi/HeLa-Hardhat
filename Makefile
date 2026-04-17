include .env
export

migrate-up:
	migrate -path internal/migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path internal/migrations -database "$(DATABASE_URL)" down 1

# make: make migrate-down
# make migrate-up

mock-service:
	mockery --name Service --dir internal/service/bank --output internal/mocks --outpkg mocks

mock-repo:
	mockery --name Repository --dir internal/repository/bank --output internal/mocks --outpkg mocks

mock-gateway:
	mockery --name Gateway --dir internal/blockchain --output internal/mocks --outpkg mocks


# make mock-service
# make mock-repo
# make mock-gateway

mock-clean-svc:
	rm -f internal/mocks/Service.go
# make mock-clean-svc

mock-clean-repo:
	rm -f internal/mocks/Repository.go

mock-clean-gateway:
	rm -f internal/mocks/Gateway.go