# システム構成図

このドキュメントは、Todoアプリケーションのシステム全体を構成するコンポーネントと、それらの間の基本的なやり取りを図で示します。

## 図

```mermaid
graph TD
    subgraph "Client Side"
        Client[Client<br>(Web Browser, etc.)]
    end

    subgraph "Server Side"
        APIServer[API Server<br>(Go, Gin)]
        Database[PostgreSQL DB]
    end

    Client -- HTTPS Request --> APIServer
    APIServer -- SQL Query --> Database
    Database -- SQL Result --> APIServer
    APIServer -- JSON Response --> Client
```

## コンポーネントの説明

- **Client (Web Browser, etc.)**: ユーザーが直接操作するインターフェース。フロントエンドアプリケーションが動作する環境（例: ウェブブラウザ）。ユーザーの操作に応じてAPIサーバーにリクエストを送信する。

- **API Server (Go, Gin)**: バックエンドの中核。Go言語とGinフレームワークで構築される。クライアントからのリクエストを受け取り、ビジネスロジックを実行し、必要に応じてデータベースと通信する。結果はJSON形式でクライアントに返却する。

- **PostgreSQL DB**: アプリケーションのデータを永続的に保存するデータベース。Todoアイテムなどの情報が格納される。APIサーバーからのクエリに応じてデータの読み書きを行う。
