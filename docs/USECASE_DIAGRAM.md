# ユースケース図

このドキュメントは、Todo APIシステムのユースケースを図で示します。

## 図

```mermaid
usecaseDiagram
    actor User
    rectangle "Todo API System" {
        User -- (Create Todo)
        User -- (Read Todos)
        User -- (Update Todo)
        User -- (Delete Todo)
    }
```

## 各ユースケースの説明

- **User**: このシステムのTodoリストを管理する利用者。
- **Create Todo**: 新しいTodo項目をシステムに登録する。
- **Read Todos**: 既存のTodo項目のリスト、または特定のTodo項目を閲覧する。
- **Update Todo**: 既存のTodo項目の内容やステータス（完了/未完了）を変更する。
- **Delete Todo**: 不要になったTodo項目をシステムから削除する。
