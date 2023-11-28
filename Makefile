include .env

migration:
	cd ./internal/database/sql/schema && goose postgres ${DB_URL} up

migration_down:
	cd ./internal/database/sql/schema && goose postgres ${DB_URL} down

sqlc:
	sqlc generate
