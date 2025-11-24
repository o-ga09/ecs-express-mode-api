.PHONY: lint
lint:
	go vet ./...
	golangci-lint run

.PHONY: test
test:
	go test -race -parallel 1 ./...

.PHONY: test-coverage
test-coverage:
	go test ./... -coverprofile=coverage.out

.PHONY: generate
generate:
	go generate ./...

# DATABASE_URL is expected to be set in the environment
ifndef DATABASE_URL
$(warning DATABASE_URL is not set. Using default local development settings)
DATABASE_URL ?= user:P@ssw0rd@tcp(127.0.0.1:3306)/develop_tavinikkiy?parseTime=true
export DATABASE_URL
endif

.PHONY: migrate-up
migrate-up:
	$(eval MIGRATE_CMD := go run cmd/migration/main.go)
	$(MIGRATE_CMD) -command up

.PHONY: migrate-down
migrate-down:
	$(eval MIGRATE_CMD := go run cmd/migration/main.go)
	$(MIGRATE_CMD) -command down

.PHONY: migrate-status
migrate-status:
	$(eval MIGRATE_CMD := go run cmd/migration/main.go)
	$(MIGRATE_CMD) -command status

.PHONY: migrate-new
migrate-new:
	$(eval MIGRATE_CMD := go run cmd/migration/main.go)
	$(MIGRATE_CMD) -command new -name $(name)

.PHONY: seed
seed:
	$(eval MIGRATE_CMD := go run cmd/migration/main.go)
	$(MIGRATE_CMD) -command seed

.PHONY: build
build:
	go build -o bin/api cmd/api/main.go

.PHONY: up
up:
	docker compose up -d db localstack

.PHONY: down
down:
	docker compose down -v

.PHONY: stop
stop:
	docker compose down

.PHONY: start
start:
	docker compose up -d
