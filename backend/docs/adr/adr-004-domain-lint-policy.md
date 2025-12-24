# ADR-XXX: Domain Boundary and Lint Policy

## Status

Accepted

## Context

本プロジェクトでは DDD およびクリーンアーキテクチャを採用している。
これまで設計ルールは暗黙知やレビューに依存しており、
集約境界やレイヤー責務の逸脱を防ぎきれなかった。

## Decision

### Domain Layer

- 純粋なビジネスルールのみを保持する
- 以下を禁止する
  - context.Context
  - time / uuid / rand
  - logger / fmt.Print
  - DB / HTTP / I/O
- 他集約への直接 import を禁止する
- 集約間連携は domain event のみ許可する

### Usecase Layer

- アプリケーション固有の処理フローを担当
- 複数集約の repository / service を組み合わせる
- valueobject の直接参照は禁止する

### Infrastructure Layer

- 技術的詳細（DB, Logger, 外部API）を担当
- domain / usecase の interface を実装する

### Import Rules

- entity / valueobject / repository は必ず alias を付与
- alias は `{aggregate}{Layer}` 形式とする
  - 例: todoVO, userRepo

### Tooling

- golangci-lint + forbidigo により設計ルールを強制
- CI で violation を error として扱う

## Consequences

- 設計違反が静的解析で検知される
- 集約の独立性が維持される
- 将来の分割・非同期化が容易になる
