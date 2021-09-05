# Synchronous Controller

リアルタイムなスライドを実現するためのめちゃくちゃすごい（自分調べ）Dapr appです。

## Env

```env
KEY= # Google IAM json
PUBSUB_PROJECT_ID=

DB_USER=
DB_NAME=
DB_PASSWORD=
```

## Tests

- 必要なもの
  - postgreSQL

```bash
# postgres start (homebrew onlu)
./tools/start_local_db.sh

# create db
create hello-slide-test -O postgres

# run test
./tools/db_test.sh [db username]

# delete db
dropdb hello-slide-test

# stop postgres
./tools/stop_local_db.sh
```

## LICENSE

[MIT](./LICENSE)
