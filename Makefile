include .env

migration:
	cd ./internal/database/sql/schema && goose postgres ${DB_URI} up

migration_down:
	cd ./internal/database/sql/schema && goose postgres ${DB_URI} down

sqlc:
	sqlc generate
