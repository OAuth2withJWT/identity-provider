ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: build
build:
	go build -o ./build/server cmd/server/main.go

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	rm -rf ./build

.PHONY: run
run:
	go run cmd/server/main.go

.PHONY: migration_up
migration_up:
	migrate -source file://db/migrations -database "postgres://$(DATABASE_USER):$(DATABASE_PASSWORD)@localhost:5432/$(DATABASE_NAME)?sslmode=disable" up $(if $(N),$(N),)

.PHONY: migration_down
migration_down:
	migrate -source file://db/migrations -database "postgres://$(DATABASE_USER):$(DATABASE_PASSWORD)@localhost:5432/$(DATABASE_NAME)?sslmode=disable" down $(if $(N),$(N),)

.PHONY: create_migration
create_migration:
	migrate create -ext sql -dir /db/migrations  $(NAME)

	