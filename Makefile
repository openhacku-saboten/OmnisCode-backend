ENV_FILE := .env
ENV := $(shell cat $(ENV_FILE))

ENV_TEST_FILE := .env.test
ENV_TEST := $(shell cat $(ENV_TEST_FILE))

.PHONY:run
run:
	$(ENV) go run main.go

.PHONY:test
test:
	$(ENV_TEST) go test -count=1 ./...

.PHONY:test-with-coverage
test-with-coverage:
	$(ENV_TEST) go test -covermode=atomic -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o ./cover.html

.PHONY:gen
gen:
	go generate ./...

.PHONY:docker-up
docker-up:
	docker-compose -f docker/docker-compose.dev.yml up --build

.PHONY:docker-down
docker-down:
	docker-compose -f docker/docker-compose.dev.yml down

.PHONY:local-db-up
local-db-up:
	docker-compose -f docker/docker-compose.test.yml up --build -d

.PHONY:migrate-up
migrate-up:
	$(ENV_TEST) sql-migrate up
