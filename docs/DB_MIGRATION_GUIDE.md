# DBマイグレーション手順書

## 1. 概要

このドキュメントは、データベーススキーマの変更を管理するための手順を定める。
本プロジェクトでは、スキーマ管理ツールとして [Atlas](https://atlasgo.io/) を使用する。

スキーマの変更は、バージョン管理されたSQLファイル（マイグレーションファイル）を通して行われる。これにより、誰がいつどのような変更を行ったかの追跡が容易になり、チーム内でのスキーマ共有や、本番環境への安全な適用が可能になる。

## 2. 使用ツール

- **Atlas**: データベーススキーマ管理ツール。
- **マイグレーションファイルディレクトリ**: `backend/migrations`

## 3. 基本的なワークフロー

1.  `backend/schema.sql` を編集して、理想のスキーマ状態を定義する。
2.  Atlasを使い、現在のデータベースの状態と `schema.sql` との差分から、マイグレーション用のSQLファイルを自動生成する。
3.  生成されたSQLファイルを確認・修正する。
4.  マイグレーションファイルをデータベースに適用する。

## 4. 新規マイグレーションファイルの作成手順

`schema.sql` を変更後、以下のコマンドを実行して差分を検出し、新しいマイグレーションファイルを生成する。

### 4.1. コマンド

```sh
# backendディレクトリで実行
atlas migrate diff [migration_name] \
  --dir "file://migrations" \
  --to "file://schema.sql" \
  --dev-url "docker://postgres/15/dev"
```

### 4.2. コマンド解説

- **`[migration_name]`**:
  マイグレーションの内容を表す名前（例: `add_user_table`, `add_priority_to_todos`）。ファイル名の一部になる。
- **`--dir "file://migrations"`**:
  マイグレーションファイルを管理するディレクトリを指定。
- **`--to "file://schema.sql"`**:
  変更後のスキーマ（あるべき姿）が定義されたファイルを指定。
- **`--dev-url "docker://postgres/15/dev"`**:
  差分計算のために一時的に起動する開発用データベースを指定。AtlasはDocker上で一時的なDBを起動し、現在のマイグレーション履歴を適用した状態を作り出し、それと `schema.sql` との差分を計算する。

### 4.3. 実行後

コマンドが成功すると、`backend/migrations` ディレクトリに `YYYYMMDDHHMMSS_[migration_name].sql` という形式で新しいファイルが作成される。
自動生成されたSQLが意図通りか必ず確認し、必要であれば手動で修正する。

## 5. マイグレーションの適用手順

保留中の（まだ適用されていない）マイグレーションファイルをデータベースに適用する。

### 5.1. ローカル開発環境への適用

`docker-compose` で起動しているローカルのPostgreSQLデータベースに適用する場合。

```sh
# backendディレクトリで実行
atlas migrate apply \
  --dir "file://migrations" \
  --url "postgres://user:password@localhost:5432/dbname?sslmode=disable"
```
※ `user`, `password`, `localhost:5432`, `dbname` は `docker-compose.yml` の設定に合わせる。

### 5.2. 本番環境への適用

**注意: 本番環境への適用は、事前にステージング環境でテストを行い、アプリケーションをメンテナンスモードにするなど、慎重に行うこと。**

```sh
# 本番DBの接続情報を使って実行
atlas migrate apply \
  --dir "file://migrations" \
  --url "[PRODUCTION_DB_URL]"
```

## 6. スキーマの状態確認

`atlas migrate status` コマンドで、データベースの現在のマイグレーション適用状況と、保留中のマイグレーションファイルを確認できる。

```sh
# ローカルDBの状態を確認する場合
atlas migrate status \
  --dir "file://migrations" \
  --url "postgres://user:password@localhost:5432/dbname?sslmode=disable"
```
