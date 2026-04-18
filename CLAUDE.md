# CLAUDE.md

Claude Code がこのリポジトリで作業する際のガイドです。

## プロジェクト概要

Go 製の Todo 管理 REST API。Clean Architecture で構成されたバックエンドサービス。

- **言語**: Go 1.25
- **フレームワーク**: Gin (HTTP), Bob (ORM), Wire + Kessoku (DI)
- **DB**: PostgreSQL / Atlas (マイグレーション)
- **ツール管理**: Aqua
- **観測性**: OpenTelemetry, zerolog

## ディレクトリ構成

```
backend/
├── cmd/server/          # エントリポイント・DI設定
├── internal/
│   ├── domain/          # エンティティ・リポジトリインターフェース
│   │   ├── todo/
│   │   ├── user/
│   │   └── apikey/
│   ├── usecase/         # ユースケース層
│   ├── infrastructure/  # DB・キャッシュなど外部依存の実装
│   │   └── rdb/         # Bob ORM モデル・Atlas マイグレーション
│   └── presentation/    # HTTP ハンドラ
└── docs/swagger/        # Swagger 生成ドキュメント
```

## よく使うコマンド

```bash
cd backend

make build          # バイナリビルド
make run            # サーバ起動
make test           # テスト実行
make lint           # golangci-lint
make air            # ホットリロード起動（開発時推奨）
make up             # Docker Compose でDB起動
make down           # DB停止
```

### コード生成

```bash
make gen-all        # 全ジェネレータ実行（Wire + Kessoku + Bob）
make gen-wire       # Wire DI コード生成
make gen-kessoku    # Kessoku DI コード生成（go generate）
make gen-bob        # Bob ORM モデル生成（bobgen-psql）
make gen-swagger    # Swagger ドキュメント生成
```

### DB マイグレーション

```bash
make migrate-diff name=<name>   # マイグレーションファイル作成
make migrate-apply              # マイグレーション適用
make migrate-hash               # ハッシュ更新
make schema-clean               # スキーマ全削除（破壊的操作）
```

## 開発ワークフロー

1. `make up` でDB起動
2. `make migrate-apply` でマイグレーション適用
3. `make air` でホットリロード起動
4. スキーマ変更時: `migrate-diff` → `migrate-apply` → `gen-bob` → `gen-all`
5. DI 変更時: `gen-all` でコード再生成

## アーキテクチャのルール

- **依存方向**: `presentation → usecase → domain ← infrastructure`
- `domain` 層は外部パッケージに依存しない
- DB アクセスは Bob ORM を通じて `infrastructure/rdb` で行う
- DI は Wire/Kessoku で管理し、手動で依存を渡さない
- API キーは `X-API-Key` ヘッダで認証

## テスト

```bash
make test           # 全テスト実行（約5〜6秒）
```

### テスト方針

層ごとにテスト戦略を分けている。

| 層 | 方式 | DB起動 |
|---|---|---|
| ドメイン（値オブジェクト・エンティティ） | ユニットテスト | 不要 |
| ユースケース | mockery 生成モックを使用 | 不要 |
| リポジトリ | testcontainers で自動起動（`TestMain` で1回） | 自動（Docker必須） |

- ドメイン・ユースケーステストは `make up` 不要でどこでも実行可能
- リポジトリテストは Docker が起動していれば自動でコンテナを起動・終了する
- `make up` はアプリ開発時のみ必要（テストでは不要）

### モック生成

```bash
mockery             # internal/domain 配下のインターフェースからモック生成
```

- 設定: `.mockery.yml`（`package mocks`、非 `_test.go` ファイルで他パッケージからインポート可能）
- 生成先: `internal/domain/**/mocks/mock_*.go`
- インターフェース変更時は `mockery` を再実行する

### テストヘルパー

- `internal/testhelper/db.go`: testcontainers を使った DB 起動・マイグレーション適用・TRUNCATE ヘルパー
- `internal/config/clock/mock_clock.go`: 固定時刻を返す `MockClock`

## Claude Code ワークフロー指針

- 複雑なタスクは `/plan` で設計してから実装する
- コンテキストが 50% 前後で `/compact` を実行する
- 複数ステップの実装はタスクリストで進捗管理する
- 問題の診断には `/doctor` を使う
- コミットはファイル単位で分けて作成する（複数ファイルをまとめない）
