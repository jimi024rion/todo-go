# APIサーバー 詳細設計書

## 1. 概要

本書は、Todo APIサーバーの内部実装に関する詳細設計を定義する。基本設計書に基づき、各レイヤーのコンポーネント、インターフェース、データ構造、および主要なロジックについて詳細に記述する。

## 2. 設計方針

本システムは、関心の分離とテスト容易性を高めるため、クリーンアーキテクチャの原則に準拠する。依存関係の方向は、 `presentation` -> `usecase` -> `domain` とし、`infrastructure` は `domain` で定義されたインターフェースを実装する。

- **`domain`**: ビジネスの核となるルールとデータ構造（エンティティ、リポジトリインターフェース）を定義。
- **`usecase`**: アプリケーション固有のビジネスロジックを実装。`domain`のエンティティとリポジトリインターフェースに依存。
- **`infrastructure`**: データベースアクセスなど、外部システムとの連携を実装。`usecase`層にDIされる。
- **`presentation`**: HTTPリクエストの受付、レスポンスの返却を担当。`usecase`層に処理を委譲。

## 3. `domain`層 詳細設計

### 3.1. エンティティ (`internal/domain/todo/entity/todo.go`)

`Todo`エンティティはビジネスドメインの中心となる。

```go
package entity

import "time"

type Todo struct {
    ID          int64
    Title       string
    Description string
    Completed   bool
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

### 3.2. リポジトリインターフェース (`internal/domain/todo/repository.go`)

`usecase`層がデータ永続化層（`infrastructure`）を抽象的に扱うためのインターフェース。

```go
package todo

import (
    "context"
    "your_project/internal/domain/todo/entity"
)

type TodoRepository interface {
    FindAll(ctx context.Context) ([]*entity.Todo, error)
    FindByID(ctx context.Context, id int64) (*entity.Todo, error)
    Create(ctx context.Context, todo *entity.Todo) (*entity.Todo, error)
    Update(ctx context.Context, todo *entity.Todo) error
    Delete(ctx context.Context, id int64) error
}
```

## 4. `usecase`層 詳細設計 (`internal/usecase/todo/`)

### 4.1. Usecaseインターフェース

`presentation`層が利用するビジネスロジックのインターフェース。

```go
package todo

import (
    "context"
    "your_project/internal/domain/todo/entity"
)

type TodoUsecase interface {
    GetAllTodos(ctx context.Context) ([]*entity.Todo, error)
    GetTodoByID(ctx context.Context, id int64) (*entity.Todo, error)
    CreateTodo(ctx context.Context, title, description string) (*entity.Todo, error)
    UpdateTodo(ctx context.Context, id int64, title, description string, completed bool) error
    DeleteTodo(ctx context.Context, id int64) error
}
```

### 4.2. `CreateTodo`メソッド処理フロー

1. `title`のバリデーション（空でないこと）を行う。
2. `entity.Todo`を生成する (`Completed`は`false`、`CreatedAt`と`UpdatedAt`は現在時刻）。
3. `TodoRepository.Create`を呼び出し、永続化する。
4. 成功すれば、作成された`Todo`エンティティを返す。

## 5. `presentation`層 詳細設計 (`internal/presentation/handler/`)

### 5.1. `TodoHandler` 構造体

`TodoUsecase`への依存をDI（Dependency Injection）で解決する。

```go
package handler

import "your_project/internal/usecase/todo"

type TodoHandler struct {
    todoUsecase todo.TodoUsecase
}

func NewTodoHandler(tu todo.TodoUsecase) *TodoHandler {
    return &TodoHandler{todoUsecase: tu}
}
```

### 5.2. `Create`メソッド 処理フロー（`POST /todos`）

1. `gin.Context`からリクエストボディ(JSON)をDTO(Data Transfer Object)にバインドする。
2. バインド時またはDTOのバリデーションでエラーがあれば `400 Bad Request` を返す。
3. `todoUsecase.CreateTodo`を呼び出す。
4. Usecaseからエラーが返却されれば、エラーの種類に応じて `500 Internal Server Error` などを返す。
5. 成功すれば、Usecaseから返された`Todo`エンティティをJSONに変換し、`201 Created`でレスポンスする。

## 6. `infrastructure`層 詳細設計 (`internal/infrastructure/repository/`)

### 6.1. `todoRepository` 構造体

`sql.DB`または`sqlx.DB`へのコネクションを持つ。

```go
package repository

import "database/sql" // or "github.com/jmoiron/sqlx"

type todoRepository struct {
    db *sql.DB
}

func NewTodoRepository(db *sql.DB) *todoRepository {
    return &todoRepository{db: db}
}
```

### 6.2. `Create`メソッド 実装方針

1. `todo`エンティティから`title`, `description`などの値を取得する。
2. `INSERT INTO todos (...) VALUES (...) RETURNING id, created_at, updated_at` のSQL文を実行する。
3. `RETURNING`で返された値をスキャンし、引数の`todo`エンティティの`ID`, `CreatedAt`, `UpdatedAt`フィールドを更新する。
4. 実行結果を返す。
