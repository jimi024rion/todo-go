# ドメインモデル図

このドキュメントは、Todo APIシステムのドメインモデルを図で示します。

## 図

ドメインモデルの中心は `Todo` エンティティです。これはアプリケーションで管理されるべき核となるデータ構造を表します。

```mermaid
classDiagram
    class Todo {
        +int id
        +string title
        +string description
        +boolean completed
        +datetime created_at
        +datetime updated_at
    }
```

## 属性の説明

- **id**: Todo項目を一意に識別するためのID。
- **title**: Todoのタイトル。
- **description**: Todoの詳細な説明。
- **completed**: Todoが完了したかどうかを示す状態。
- **created_at**: Todoが作成された日時。
- **updated_at**: Todoが最後に更新された日時。
