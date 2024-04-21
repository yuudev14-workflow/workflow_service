.PHONE: migrate

DATABASE_URL := postgres://postgres:password@localhost/workflow?sslmode=disable
TEST_DATABASE_URL := postgres://postgres:password@localhost/test_workflow?sslmode=disable
MIGRATIONS_PATH := file://db/migrations

migrate:
		migrate -source $(MIGRATIONS_PATH) -database $(DATABASE_URL) up

migrate-down:
		migrate -source $(MIGRATIONS_PATH) -database $(DATABASE_URL) down

test:
	migrate -source $(MIGRATIONS_PATH) -database $(TEST_DATABASE_URL) up && \
	go test ./...

swag:
	swag init