include .env
export

migrate-up:
	migrate -path internal/migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path internal/migrations -database "$(DATABASE_URL)" down 1
