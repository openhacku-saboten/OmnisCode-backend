#!/bin/sh

go get -v github.com/rubenv/sql-migrate/...

export $(cat .env)
while :
do
    sql-migrate up
    if [ $? -eq 0 ]; then
        break
    fi
    sleep 2
done

go run main.go
