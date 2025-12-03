# デプロイ手順書

## 1. 概要

このドキュメントは、Todo APIサーバーをローカル開発環境で実行する手順、およびGoogle Cloud Platform (GCP) 上のCloud Runへデプロイする手順を記述する。

## 2. ローカル環境での実行

### 2.1. 前提条件

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

上記がインストールされ、Dockerデーモンが実行中であること。

### 2.2. 実行手順

1. プロジェクトルートの `backend` ディレクトリに移動する。
    ```sh
    cd backend
    ```

2. `docker-compose` を使用して、APIサーバーとデータベースのコンテナをビルドし、バックグラウンドで起動する。
    ```sh
    docker-compose up --build -d
    ```

3. コンテナの起動状況を確認する。
    ```sh
    docker-compose ps
    ```
    `State`が`Up`になっていれば正常に起動している。

### 2.3. 動作確認

`curl` コマンドで `localhost:8080/api/v1/todos` (ポートは`docker-compose.yml`の設定に依存) にアクセスし、空の配列 `[]` が返ってくることを確認する。

```sh
curl http://localhost:8080/api/v1/todos
```

### 2.4. 停止手順

```sh
# `backend` ディレクトリで実行
docker-compose down
```

## 3. GCP (Cloud Run) へのデプロイ手順 (想定)

これは将来的にGCPへデプロイする際の標準的な手順を示す。

### 3.1. 前提条件

- [Google Cloud SDK (gcloud CLI)](https://cloud.google.com/sdk/docs/install)がインストール・設定済みであること。
- 課金が有効なGCPプロジェクトが作成済みであること。
- 以下のGCP APIが有効になっていること。
    - Cloud Build API
    - Artifact Registry API
    - Cloud Run API
    - Cloud SQL Admin API
    - Serverless VPC Access API

### 3.2. ステップ1: Dockerfileの作成

プロジェクトの `backend` ディレクトリに、以下の内容で `Dockerfile` を作成する。

```Dockerfile
# --- Build Stage ---
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Goモジュールの依存関係をキャッシュ
COPY go.mod go.sum ./
RUN go mod download

# アプリケーションのソースコードをコピー
COPY . .

# アプリケーションをビルド
# CGO_ENABLED=0: 静的リンクバイナリを生成
# -ldflags="-s -w": シンボルテーブルとデバッグ情報を削除し、バイナリサイズを削減
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/app

# --- Release Stage ---
FROM alpine:latest

WORKDIR /app

# ビルドステージからコンパイル済みのバイナリをコピー
COPY --from=builder /app/server /app/server

# ポート8080を公開
EXPOSE 8080

# アプリケーションを実行
CMD ["/app/server"]
```

### 3.3. ステップ2: コンテナイメージのビルドとプッシュ

Cloud Buildを使い、コンテナイメージをビルドしてArtifact Registryにプッシュする。

```sh
# 環境変数を設定
export PROJECT_ID=[YOUR_PROJECT_ID]
export REGION=[YOUR_REGION] #例: asia-northeast1
export REPO_NAME=todo-go-repo
export IMAGE_NAME=todo-api

# Artifact Registryにリポジトリを作成 (初回のみ)
gcloud artifacts repositories create ${REPO_NAME} \
    --repository-format=docker \
    --location=${REGION}

# Cloud Buildを実行 (`backend`ディレクトリで実行)
gcloud builds submit . \
    --tag ${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPO_NAME}/${IMAGE_NAME}:latest
```

### 3.4. ステップ3: Cloud SQL for PostgreSQLのセットアップ

（詳細はGCPドキュメントを参照）
1. Cloud SQLインスタンスを作成する。
2. `todos` データベースを作成する。
3. APIサーバーが使用するデータベースユーザーを作成する。

### 3.5. ステップ4: Cloud Runサービスのデプロイ

```sh
# 環境変数を設定
export SERVICE_NAME=todo-api-service
export DB_CONNECTION_STRING="host=/cloudsql/[INSTANCE_CONNECTION_NAME] user=[DB_USER] password=[DB_PASS] dbname=todos sslmode=disable"

# Cloud Runサービスをデプロイ
gcloud run deploy ${SERVICE_NAME} \
    --image=${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPO_NAME}/${IMAGE_NAME}:latest \
    --platform=managed \
    --region=${REGION} \
    --allow-unauthenticated \
    --set-env-vars="DB_DSN=${DB_CONNECTION_STRING}" \
    --vpc-connector=[YOUR_VPC_CONNECTOR_NAME] # Cloud SQL接続に必要
```
このコマンドにより、コンテナがCloud Run上で起動し、公開URLが発行される。

