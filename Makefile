#  make test ENV=".env.test"

include $(ENV)
export

.PHONY: migrate
.PHONY: migrate-down
.PHONY: test
.PHONY: swag


DATABASE_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST)/$(DB_NAME)?sslmode=disable
MIGRATIONS_PATH := file://db/migrations

migrate:
		migrate -source $(MIGRATIONS_PATH) -database $(DATABASE_URL) up

migrate-down:
		migrate -source $(MIGRATIONS_PATH) -database $(DATABASE_URL) down

test:
	migrate -source $(MIGRATIONS_PATH) -database $(DATABASE_URL) up 
	go test -coverprofile=coverage.out -v ./... 

swag:
	swag init
