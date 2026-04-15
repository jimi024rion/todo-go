# Stage 1: Build
FROM golang:1.25-alpine AS builder

WORKDIR /app

# 依存関係を先にコピーしてキャッシュを活用
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
  -ldflags="-s -w" \
  -trimpath \
  -o /server \
  ./cmd/server

# Stage 2: Runtime (distrolessで最小イメージ)
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /server /server

EXPOSE 8080

CMD ["/server"]
