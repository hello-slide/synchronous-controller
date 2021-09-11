# Synchronous Controller

リアルタイムなスライドを実現するためのめちゃくちゃすごい（自分調べ）Dapr appです。

## Env

```env
DB_USER=
DATABASE_NAME=
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

   ```jsonc
   # req
   {
       "type": "0"
   }

   # res
   {
       "type": "0",
       "version": "1.0",
       "id": "",
   }
   ```

1. セッション開始リクエスト（visitor）
   - サーバーから同じ値で返答あり

   ```jsonc
   # req
   {
       "type": "1",
       "id": ""
   }

   # res
   {
       "type": "1",
       "version": "1.0",
   }
   ```

2. 参加者数取得

    ```jsonc
    # res
    {
        "type": "2",
        "visitors": "10",
    }
    ```

3. 回答取得

    ```jsonc
    # res
    {
        "type": "3",
        "answers": [
            {
                "id": "",
                "user_id": "",
                "name": "",
                "answer": "",
            },
            {
                "id": "",
                "user_id": "",
                "name": "",
                "answer": "",
            },
            ...
        ]
    }
    ```

4. 新しいトピック割当

    ```jsonc
    # req
    {
        "type": "4",
        "topic": "",
    }
    ```

5. トピック送信（visitor）

    ```jsonc
    # res
    {
        "type": "5",
        "topic": "",
    }
    ```

6. 回答（visitor）

    ```json
    # req
    {
        "type": "6",
        "answer": "",
        "name": "",
    }
    ```

## LICENSE

[MIT](./LICENSE)
