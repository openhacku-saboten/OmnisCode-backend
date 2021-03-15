ENV_TEST_FILE := .env.test
ENV_TEST := $(shell cat $(ENV_TEST_FILE))

.PHONY:run
run:
	go run main.go

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
