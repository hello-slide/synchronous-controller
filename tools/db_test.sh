#! /bin/bash

if [ $# != 2 ]; then
    echo "example: sh ./test.sh [db username] [password]"
    exit 1
fi

export LOCAL_TEST="db"
export DB_USER=$1
export DB_PW=$2

go test -v ./...
