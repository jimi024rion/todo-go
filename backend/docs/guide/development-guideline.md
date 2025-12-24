# Development Guideline

（開発者向け実装ガイドライン / 超詳細版）

本ドキュメントは、本プロジェクトにおける **DDD（ドメイン駆動設計）** および
**クリーンアーキテクチャ** を前提とした設計・実装・運用ルールを定義する。

目的は以下である。

- アーキテクチャ境界を**人為ミスに依存せず守る**
- 設計判断のブレを抑え、**中長期での保守性・拡張性を最大化する**
- 新規参画者が**迷わず正しい実装を選択できる状態**を作る

## 1. 基本原則（必ず守る）

### 1.1 依存関係逆転の原則（DIP）

- 依存は **外 → 内** のみ許可する
- `domain` は **何にも依存しない**
- `usecase` は `domain` のみに依存する
- `presentation / infrastructure` は `usecase / domain` に依存する

> **重要**
> これらのルールは *設計規約* ではなく *破ると事故る前提条件* である。

### 1.2 レイヤー越境禁止

以下は **設計上の違反** とする。

- presentation → infrastructure の直接呼び出し
- usecase → infrastructure の実装依存
- domain → context / DB / Web / 外部API 依存

これらは **golangci-lint（depguard）で強制**される。

## 2. ディレクトリ構成と責務

```text
internal/
├── domain/
│   ├── model/
│   ├── service/
│   ├── repository/
│   └── value/
├── usecase/
│   ├── input/
│   ├── output/
│   └── interactor/
├── presentation/
│   ├── http/
│   └── middleware/
└── infrastructure/
    ├── rdb/
    ├── external/
    └── repository/
```

## 3. 各レイヤーの詳細ガイドライン

## 3.1 Domain 層

### 3.1.1 責務

- 業務ルール・制約・概念の表現
- 不変条件（Invariant）の保証
- ビジネスロジックの中心

### 3.1.2 禁止事項（絶対）

- `context.Context` の使用
- DB / SQL / ORM 依存
- HTTP / Gin / JSON 依存
- ログ出力
- 時刻の直接取得（`time.Now()`）

### 3.1.3 設計指針

#### Entity

- 同一性（ID）を持つ
- 状態変更は **必ずメソッド経由**
- Setter を生やさない

```go
func (u *User) ChangeEmail(email Email) error
```

#### Value Object

- 不変（Immutable）
- 構造体は `private field`
- バリデーションは生成時のみ

```go
func NewEmail(v string) (Email, error)
```

#### Domain Service

- Entity / VO に自然に属さないルールのみ
- ステートレス

## 3.2 Usecase 層

### 3.2.1 責務

- ユースケースの実行制御
- トランザクション境界
- Repository / 外部IFの抽象依存

### 3.2.2 設計ルール

- **1 ユースケース = 1 Interactor**
- 入出力は DTO（Input / Output）で分離
- `context.Context` は usecase からのみ扱う

```go
type CreateUserUsecase interface {
    Execute(ctx context.Context, input CreateUserInput) (CreateUserOutput, error)
}
```

### 3.2.3 禁止事項

- HTTP ステータスコードの知識
- JSON / binding / validation
- Gin 依存
- ORM の構造体露出

## 3.3 Presentation 層

### 3.3.1 責務

- HTTP リクエスト/レスポンス変換
- 認証・認可
- 入力バリデーション
- エラーハンドリング（HTTP変換）

### 3.3.2 設計指針

- Handler は **薄く**
- 業務判断は usecase に委譲
- Request / Response 専用 struct を定義

```go
func (h *UserHandler) Create(c *gin.Context)
```

### 3.3.3 禁止事項

- DB / Repository 直接操作
- ビジネスルール実装
- トランザクション制御

## 3.4 Infrastructure 層

### 3.4.1 責務

- 外部技術との接続
- Repository / Client 実装
- ORM / SDK のカプセル化

### 3.4.2 設計指針

- domain.repository の実装のみ提供
- ORM モデルは **外に漏らさない**
- mapping を明示的に書く

```go
type UserRepository struct {}
```

### 3.4.3 禁止事項

- Gin 依存
- usecase への逆依存
- ドメインロジック記述

## 4. 命名規則

### 4.1 共通

- Package 名は **単数・小文字**
- 省略語は禁止（ctx, err は除外）
- bool は `Is/Has/Can` プレフィックス

### 4.2 Interface

- 役割名 + 動詞不要
- `Repository`, `Usecase`, `Service`

```go
type UserRepository interface {}
```

## 5. エラーハンドリング

- domain: `error` を返すのみ
- usecase: エラー分類を行う
- presentation: HTTP に変換

```go
errors.Is(err, domain.ErrUserNotFound)
```

## 6. ロギング

- domain/usecase ではログ禁止
- presentation / infrastructure のみ可
- logger は DI 経由

## 7. フォーマット & Lint

- gofmt / goimports 必須
- gci による import grouping 強制
- golangci-lint v2 準拠

## 8. テスト指針

### 8.1 Domain

- ロジック中心
- Mock 不要

### 8.2 Usecase

- Repository は Mock
- 分岐網羅

### 8.3 Presentation

- Handler 単体 or E2E

## 9. よくあるアンチパターン

- usecase が肥大化
- domain が DTO 化
- infrastructure にロジック流出
- Handler が God Object

## 10. 最後に

このガイドラインは **自由を奪うためではなく、判断コストを下げるためのもの** である。
迷った場合は **「依存の向き」と「責務」** に立ち返ること。

破りたくなった場合は、まず ADR を書くこと。
