# アーキテクチャ
## 概要
フロントエンドからのHTTP(S)リクエストをバックエンドのAPIサーバが受け取り、データの永続化のためにDBサーバと通信を行います。その後、結果をユーザに返却します。

```
フロントエンド
|   ^
v   |
APIサーバ(このリポジトリ)
|   ^
v   |
DBサーバ
```

## API仕様
下記の2種類のフォーマットで提供しています。

- [Web版(GitHub Pages)](https://openhacku-saboten.github.io/OmnisCode-backend/index.html)
- [YAML版](https://github.com/openhacku-saboten/OmnisCode-backend/blob/main/docs/swagger.yaml)

## アプリケーションアーキテクチャ
クリーンアーキテクチャを採用しました。

### ディレクトリ構成
```
├── config       # 環境変数の再構築までが責務です
├── controller   # HTTPリクエストを受け取り、エンティティとしてユースケースに渡すまでが責務です
├── docker       # Dockerに関する操作の設定ファイル群が含まれます
├── docs         # ドキュメント類が含まれます
├── domain
│   ├── entity   # 対象とそれを操作するロジックであるエンティティの定義が含まれます
│   └── service  # エンティティ自体が主体となる操作ではないエンティティの操作が責務です
├── infra        # 永続化と再構築が責務です
├── log          # ロガー定義が含まれます
├── migrations   # DBのマイグレーションファイル群が含まれます
├── repository   # 依存関係の逆転のためのインタフェース定義が含まれます
└── usecase      # コントローラから受け取ったエンティティを適切なレポジトリに振り分けるまでが責務です
```

### HTTPリクエストからDBアクセスをしてレスポンスを返すまでの流れ
各層同士では、`domain/entity/`で定義されたエンティティを最小単位としてやり取りを行う。
```
main.go
|    ^
v    |
controller/*.go
|    ^
v    |
usecase/*.go
|    ^
v    |
infra/*.go
|    ^
v    |
DBサーバ
```
