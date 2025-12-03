# --- Build Stage ---
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Goモジュールの依存関係をキャッシュ
COPY go.mod go.sum ./
RUN go mod tidy

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
