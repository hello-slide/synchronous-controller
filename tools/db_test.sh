#! /bin/bash

if [ $# != 1 ]; then
    echo "example: sh ./test.sh [db username]"
    exit 1
fi

export LOCAL_TEST="db"
export DB_USER=$1

go test -v ./...
