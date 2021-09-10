#! /bin/bash

if [ $# != 1 ] && [ $# != 2 ]; then
    echo "example: sh ./test.sh [db username] ([password])"
    exit 1
fi

export LOCAL_TEST="db"
export DB_USER=$1

if [ -z "$2" ]; then
    export DB_PW=$2
fi

go test -v ./...
