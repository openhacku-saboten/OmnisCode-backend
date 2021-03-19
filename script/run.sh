#!/bin/sh
set -ue

go get -v github.com/rubenv/sql-migrate/...

export $(cat .env)
timeout 30 sh -c "until nc -vz $DB_HOST $DB_PORT; do sleep 1; done" && sql-migrate up

go run main.go
