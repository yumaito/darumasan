# Darumasan

適当に作っただるまさんがころんだサーバーです。
websocketで通信します。

### 環境
* golang 1.8
* glide
```sh
$ curl https://glide.sh/get | sh
```

### 起動方法
```sh
$ glide install
$ go run main.go
```
これで`localhost:8080`でサーバーが立ちます

### エンドポイント

* プレイヤー端末 `ws://{domain}/client`
* 鬼端末 `ws://{domain}/curator`

### メッセージの形式

* サーバーへ送るとき

  + プレイヤーはアウトかどうか
  + 鬼は判定中かどうか
```json
{"status": true}
```

* サーバーからのメッセージ

初回接続時にも返されます
* fromはメッセージのトリガーになったクライアントのID
* toは自分のID
* client_type 1はプレイヤー, 2は鬼
* clientsはと接続しているIDの配列
* dead_clientsは既にアウトになったIDの配列
* curator_idは鬼のID
* is_watchedは判定中かどうか
```json
{
    "from":"b33056f5a7",
    "to":"b33056f5a7",
    "client_type":1,
    "clients":["b33056f5a7","d210d0acec"],
    "dead_clients":["b33056f5a7"],
    "curator_id":"d210d0acec",
    "is_watched":true
}
```

プレイヤーからメッセージが送られると、送ったプレイヤーと鬼に対して上記のメッセージが送信されます。
鬼が送った場合は鬼にのみ返信されます
