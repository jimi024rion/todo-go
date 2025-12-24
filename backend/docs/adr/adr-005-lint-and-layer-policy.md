# ADR-004: Lint 方針およびレイヤー依存ルールの策定

## Status

Accepted

## Context

本プロジェクトでは、DDD およびクリーンアーキテクチャを採用している。
しかし、以下の課題があった。

- 設計ルールがレビュー依存になりやすい
- domain 層に技術的関心（context, time, logger 等）が混入する
- import が増えるにつれ、レイヤー境界が視認しづらくなる
- 人によって「やっていい・ダメ」の判断が揺れる

これらを **lint によって自動的に検知・抑止** するため、明確な方針が必要となった。

---

## Decision

### 1. golangci-lint v2 を採用する

- CI / ローカルで同一ルールを適用
- Phase3 では「新規差分のみ」を対象とする

---

### 2. レイヤー依存は depguard で強制する

#### domain 層

- 技術的関心を持たない
- context / time / uuid / logger / DB 依存を禁止
- 他レイヤーへの依存を禁止

#### usecase 層

- context を受け取る唯一の内側レイヤー
- infrastructure 実装への直接依存は禁止

#### presentation 層

- HTTP 入出力のみを責務とする
- DB / infrastructure 直接操作は禁止

#### infrastructure 層

- 技術詳細の実装のみ
- presentation には依存しない

---

### 3. 禁止 API は forbidigo で管理する

- fmt.Print / print / log.Fatal / Panic を禁止
- domain 層での time.Now / uuid.New を禁止
- 「便利だが事故る API」を明示的に排除する

---

### 4. import 整理に gci を採用する

- import の並び順を機械的に統一
- 標準 / 外部 / internal を明確に分離
- alias 命名は lint ではなくガイドラインで半強制とする

---

## Consequences

### 良い点

- 設計レビューの負荷が大幅に下がる
- domain 層の純度が保たれる
- 新規参加者がルールを「読まなくても守れる」

### 悪い点・制約

- 初期導入時は lint エラーが多く出る
- 一部ルールは思想レベルに留まる（alias 命名など）

---

## Notes

本 ADR は「設計の自由を奪う」ためではなく、
**設計判断を人からツールに移譲する** ことを目的とする。

設計議論は「業務」に集中させるべきであり、
単純な禁止事項は lint によって自動化する。
