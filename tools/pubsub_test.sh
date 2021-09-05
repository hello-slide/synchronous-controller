#! /bin/bash

if [ $# != 1 ]; then
    echo "example: sh ./test.sh [google iam json path]"
    exit 1
fi

export LOCAL_TEST="pubsub"
export IAM_PATH=$1

go test -v ./...
