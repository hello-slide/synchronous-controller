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

## Websocket types

0. セッション開始リクエスト（Host）
   - サーバーから同じ値で返答あり
1. セッション開始リクエスト（visitor）
   - サーバーから同じ値で返答あり
2. 参加者数取得
3. 回答取得
4. 新しいトピック割当
5. トピック送信（visitor）
6. 回答（visitor）

## LICENSE

[MIT](./LICENSE)
