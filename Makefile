

.PHONY:run
run:
	go run main.go

.PHONY:test
test:
	go test -count=1 ./...

.PHONY:test-with-coverage
test-with-coverage:
	go test -covermode=atomic -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o ./cover.html

.PHONY:gen
gen:
	go generate ./...
