# 設計者向けガイドライン（DDD / クリーンアーキテクチャ）

## 1. 本ガイドラインの目的

本ガイドラインは以下を目的とする。

- ドメインモデルの純度を長期的に保つ
- レイヤー責務の混濁を防ぐ
- 実装者ごとの判断ブレを最小化する
- レビューを「設計」ではなく「業務」に集中させる

本プロジェクトでは **lint を「設計ルールの自動番人」** と位置付ける。

---

## 2. 採用アーキテクチャ概要

### 採用方針

- ドメイン駆動設計（DDD）
- クリーンアーキテクチャ
- レイヤー分割は「技術」ではなく「責務」で行う

```

presentation → usecase → domain
↑
infrastructure

````

### 依存方向の原則

- **内側は外側を知らない**
- 例外は「インターフェース経由のみ」

---

## 3. レイヤー別 責務と禁止事項

---

## 3.1 domain 層

### 役割

- 業務ルールの表現
- 不変条件（Invariant）の保持
- 業務用語の集約（ユビキタス言語）

### 含めてよいもの

- Entity
- Value Object
- Domain Service
- Repository Interface
- 業務例外（Error）

### 明示的に禁止するもの

| 分類 | 理由 |
|---|---|
| context.Context | 技術的関心であり業務ではない |
| time.Now / uuid | 非決定性がテストと理解を壊す |
| logger | 副作用は業務ではない |
| DB / SQL | 永続化は業務ではない |
| Gin / HTTP | 入出力形式は業務ではない |

### 設計ルール

- **domain は「純粋な Go」**
- 副作用を持たない
- 状態遷移は必ずメソッド経由

```go
// OK
func (t *Todo) Complete() error {
    if t.status == Done {
        return ErrAlreadyDone
    }
    t.status = Done
    return nil
}
````

---

## 3.2 usecase 層

### 役割

- アプリケーション固有の処理フロー
- トランザクション境界
- 権限・制約の調整

### 含めてよいもの

- context.Context
- Repository Interface への依存
- 外部サービス Interface

### 禁止事項

| 禁止            | 理由                      |
| ------------- | ----------------------- |
| Gin 直接操作      | HTTP は presentation の責務 |
| DB 実装         | 永続化手段に依存しないため           |
| domain ルールの重複 | 業務は domain に閉じ込める       |

### 設計ルール

- **usecase = 手続きのオーケストレーター**
- domain ロジックを書かない
- 例外は domain error を返す

```go
func (uc *CreateTodoUsecase) Execute(
    ctx context.Context,
    input CreateTodoInput,
) error {
    todo, err := domain.NewTodo(input.Title)
    if err != nil {
        return err
    }
    return uc.repo.Save(ctx, todo)
}
```

---

## 3.3 presentation 層

### 役割

- HTTP / JSON / Path / Query の変換
- 入出力の validation
- usecase 呼び出し

### 禁止事項

| 禁止           | 理由            |
| ------------ | ------------- |
| domain ロジック  | 責務外           |
| DB 操作        | infra の責務     |
| context の再生成 | Gin の ctx を使う |

### 設計ルール

- handler は薄く保つ
- domain 型を直接 JSON にしない

```go
func (h *TodoHandler) Create(c *gin.Context) {
    var req CreateTodoRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        respondBadRequest(c, err)
        return
    }

    err := h.uc.Execute(c.Request.Context(), req.toInput())
    if err != nil {
        respondError(c, err)
        return
    }

    c.Status(http.StatusCreated)
}
```

---

## 3.4 infrastructure 層

### 役割

- 技術詳細の実装
- DB / 外部 API / Logger / Cache

### 設計ルール

- interface を満たすだけ
- domain を汚染しない
- エラーは domain error に変換する

```go
func (r *TodoRepository) Save(
    ctx context.Context,
    todo *domain.Todo,
) error {
    // sqlboiler を使った永続化
}
```

---

## 4. import 命名ルール（半強制）

### 基本方針

- import は「可読性のために明示する」
- lint ではなくガイドラインで制御

### 推奨命名

| 種類            | 命名例          |
| ------------- | ------------ |
| domain entity | `todoEntity` |
| value object  | `todoVO`     |
| usecase       | `todoUC`     |
| repository    | `todoRepo`   |

```go
import (
    todoEntity "internal/domain/todo/model/entity"
    todoVO     "internal/domain/todo/model/valueobject"
    todoUC     "internal/usecase/todo"
)
```

---

## 5. lint の位置付け

### lint で強制するもの

- レイヤー依存（depguard）
- 禁止 API（forbidigo）
- 複雑度・肥大化（gocyclo / funlen）
- セキュリティ事故（gosec）

### lint で強制しないもの

- alias 命名
- domain 内部設計の美学
- 集約設計の粒度

👉 **lint は設計を「補助」するもの**

---

## 6. 設計判断に迷ったときの原則

1. それは業務か？技術か？
2. 将来差し替わる可能性はあるか？
3. domain に置いた場合、テストは楽か？
4. 5 年後に読んで意味が通じるか？

---

## 7. 最後に

このガイドラインは「縛る」ためのものではない。
**設計の自由を守るために、不要な自由を減らす** ためのものである。

ルールに例外が必要な場合は、必ず ADR に残すこと。
