# シーケンス図

このドキュメントは、主要なユースケースにおけるコンポーネント間のインタラクションをシーケンス図で示します。

## 1. Todoの新規作成 (`POST /todos`)

この図は、クライアントが新しいTodoを作成するリクエストを送った際の、システム内部の処理の流れを示します。

```mermaid
sequenceDiagram
    participant Client
    participant Router as Router (Gin)
    participant Handler
    participant Usecase
    participant Repository
    participant Database as DB (PostgreSQL)

    Client->>+Router: POST /api/v1/todos (JSON)
    Router->>+Handler: HandleCreateTodo(context)
    Handler->>Handler: Bind & Validate Request
    Handler->>+Usecase: CreateTodo(title, description)
    Usecase->>Usecase: Create Todo Entity
    Usecase->>+Repository: Save(todo)
    Repository->>+Database: INSERT INTO todos ...
    Database-->>-Repository: (Success, return ID)
    Repository-->>-Usecase: (Created Todo Entity)
    Usecase-->>-Handler: (Created Todo Entity)
    Handler-->>-Router: 201 Created (JSON Response)
    Router-->>-Client: 201 Created (JSON Response)
```

## 2. Todoの一覧取得 (`GET /todos`)

この図は、クライアントがTodoリスト全体を要求した際の、システム内部の処理の流れを示します。

```mermaid
sequenceDiagram
    participant Client
    participant Router as Router (Gin)
    participant Handler
    participant Usecase
    participant Repository
    participant Database as DB (PostgreSQL)

    Client->>+Router: GET /api/v1/todos
    Router->>+Handler: HandleGetTodos(context)
    Handler->>+Usecase: FindAllTodos()
    Usecase->>+Repository: FindAll()
    Repository->>+Database: SELECT * FROM todos
    Database-->>-Repository: (Todo Records)
    Repository-->>Repository: Map Records to Entities
    Repository-->>-Usecase: (List of Todo Entities)
    Usecase-->>-Handler: (List of Todo Entities)
    Handler-->>-Router: 200 OK (JSON Response)
    Router-->>-Client: 200 OK (JSON Response)
```
