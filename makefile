DB_URL=postgres://emchepe:welcome1@localhost:5432/test_user_management?sslmode=disable

ifeq ($(OS),Windows_NT)
	DETECTED_OS := Windows
else
	DETECTED_OS := $(shell sh -c 'uname 2>/dev/null || echo Unknown')
endif

.SILENT: help
help:
	@echo "Usage: make [command]"
	@echo ""
	@echo "Commands:"
	@echo " migration-create name={name}  Create migration"
	@echo " migration-up                  Up migrations"
	@echo " migration-down                Down last migration"
	@echo " go-run                        Run Project"
	@echo " docker-up                     Run with docker"
	@echo " docker-down                   Stop docker"

# Build

.SILENT: migration-create
migration-create:
	@migrate create -ext sql -dir ./migrations -seq $(name)

# Up migration

.SILENT: migration-up
migration-up:
	@migrate -database $(DB_URL) -path ./migrations up

# Down migration

.SILENT: migration-down
migration-down:
	@migrate -database $(DB_URL) -path ./migrations down 1

.SILENT: go-run
go-run:
	@go run main.go

.SILENT : docker-up
docker-up:
	@docker compose up -d

.SILENT : docker-down
docker-down:
	@docker compose down -v

.DEFAULT_GOAL := help
