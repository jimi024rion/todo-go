# `zerolog` を活用したカスタムスタックトレース実装

## 1. 概要

本プロジェクトでは、構造化ロガーとして `zerolog` を採用しています。エラー発生時に、デバッグを容易にするための詳細なスタックトレースをログに出力する必要があります。

従来、Goコミュニティでは `pkg/errors` ライブラリがこの目的で広く使われていましたが、現在はアーカイブされ、開発が終了しています。

そこで、本プロジェクトでは `pkg/errors` に依存せず、Goの標準ライブラリ `runtime` と `zerolog` の拡張機能 `ErrorStackMarshaler` を活用し、独自のスタックトレース実装を構築しました。

## 2. アーキテクチャ

本実装のアーキテクチャは、**「エラー生成時の情報記録」** と **「ログ出力時の情報整形」** という2つのフェーズに明確に分離されています。

![Architecture Diagram](https://mermaid.ink/svg/eyJjb2RlIjoiZ3JhcGggVERcbiAgICBzdWJncmFwaCBcIkVscm9yIENyZWF0aW9uIFBoYXNlXCJcbiAgICAgICAgQVtFcnJvciBPY2N1cnNcbiAgICAgICAgQVxuICAgICAgICBcdC0tPiB8XCJlcnJzLk5ld0VyclwiIGNhbGxlZHwgQihgZXJyczo6TmV3RXJyYClcbiAgICAgICAgQlxuICAgICAgICBcdC0tPiB8XCJydW50aW1lLkNhbGxlcnMoMikgZXRjLlwiIHwgQ3tbXCJFYXJyYCBzdHJ1Y3Qgd2l0aFxuICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgYGNhbGxlcnMgW111aW50cHJgXCJdXG4gICAgZW5kXG5cbiAgICBzdWJncmFwaCBcIkxvZ2dpbmcgUGhhc2VcIlxuICAgICAgICBEXFtsb2dnZXIuRXJyb3JMb2coZXJyKVxuICAgICAgICAgICAgICAgICAgICAgfCBsb2cuU3RhY2soKV0gLS0-IEV7XCJ6ZXJvbG9nLkVycm9yU3RhY2tNYXJzaGFsZXJcbiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICBzcGVjaWZpZWQgYXMgYGVycnM6Ok1hcnNoYWxTdGFja2BcIl19XG4gICAgICAgIEVcbiAgICAgICAgXHQtLT4gfFwiUnVudGltZSBMaWJyYXJ5IFVzZWQgZm9yIEZvcm1hdHRpbmdcIiB8IEZbXCJGcmFtZXMgKGVnLiwgZnVuY05hbWUsIGxpbmUpXCJdXG4gICAgICAgIEZcbiAgICAgICAgXHQtLT4gRyhbXCJGcmFtZXMgYXJlIGZvcm1hdHRlZCBcbiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICBpbnRvIGEgaHVtYW4tcmVhZGFibGUgZm9ybWF0XCJdKVxuICAgIGVuZFxuXG4gICAgQyAtLT4gfFwiRXJyIG9iamVjdCBwYXNzZWRcIiB8IERcblxuICAgIHN0eWxlIEEgZmlsbDojZmNjLCBzdHJva2U6IzMzMywgc3Ryb2tlLXdpZHRoOjJweFxuICAgIHN0eWxlIEcgZmlsbDojY2ZjLCBzdHJva2U6IzMzMywgc3Ryb2tlLXdpZHRoOjJweFxuXG4iLCJtZXJtYWlkIjp7InRoZW1lIjoiZGVmYXVsdCJ9fQ)

1. **エラー生成時**: `errs.NewErr` が呼ばれると、`runtime.Callers` を使って現在のコールスタックの「生データ」(`[]uintptr`) を取得し、`errs.Err` 構造体に保持させます。
2. **ログ出力時**: `logger.ErrorLog(err)` の中で `.Stack()` が呼ばれると、`zerolog` はグローバルに設定された `ErrorStackMarshaler` を実行します。本プロジェクトでは、これを `errs.MarshalStack` に設定しています。
3. `errs.MarshalStack` は、渡された `err` オブジェクトから生のスタックデータ (`[]uintptr`) を取り出し、`runtime.CallersFrames` などを使って、関数名・ファイル名・行番号といった人間が読める形式に整形して返します。

## 3. 実装の詳細

### `errs.go`

このファイルは、スタックトレース付きエラーの定義と、そのスタックを整形するロジックの責務を持ちます。

#### `type Err struct { ... }`

カスタムエラーの型定義です。`callers []uintptr` フィールドがスタックトレースの「生データ」を保持します。

#### `func NewErr(...) error`

エラーを生成する関数です。`runtime.Callers(2, ...)` を呼び出し、`NewErr` 自身と `runtime.Callers` を除いた、呼び出し元のスタック情報を取得して `Err` 構造体に格納します。

#### `func MarshalStack(...) any`

`zerolog` から呼び出されるスタック整形用の関数です。以下の処理を行います。

1. `errors.As` で渡された `error` から `*Err` 型を取り出します。
2. `runtime.CallersFrames` で生のスタックデータ (`[]uintptr`) を、詳細な `runtime.Frame` に変換するためのイテレータを取得します。
3. ループで `Frame` を一つずつ取り出し、`formatFuncName` で関数名を整形し、`filepath.Base` でファイル名を抽出し、最終的に `[]map[string]any` の形式に変換して返します。

#### `func formatFuncName(...) string`

`runtime.Frame.Function` が返すフルパスの関数名（例: `github.com/foo/bar.MyFunc`）を、可読性の高い短い形式（例: `MyFunc`）に整形するヘルパー関数です。

### `logger.go`

#### `func InitializeLogger()`

アプリケーション起動時に一度だけ呼ばれ、`zerolog` のグローバル設定を行います。ここで `zerolog.ErrorStackMarshaler = errs.MarshalStack` を設定することで、スタックトレースの整形処理を `errs.MarshalStack` に委譲しています。

## 4. 他の選択肢との比較

### vs `pkg/errors` (アーカイブ済み)

- **優れている点**:
  - **依存関係の排除**: 開発が終了したライブラリへの依存がなくなりました。これは長期的なメンテナンス性において大きなメリットです。
  - **透明性と制御性**: スタック取得・整形のロジックが全てプロジェクト内のコードとなり、挙動が明確でカスタマイズも容易です。
- **劣っている点**:
  - **規約への依存**: エラー生成時に常に `errs.NewErr` を使う規約を開発者が守る必要があります。標準の `errors.New` 等ではスタックトレースが付与されません。
  - **機能のシンプルさ**: `pkg/errors` が提供していた `Wrapf` や多様なフォーマット指定 (`%+v` など) といった便利なヘルパー機能はありません。

### vs `github.com/cockroachdb/errors`

- **優れている点 (自前実装が)**:
  - **軽量・シンプル**: 我々の実装は、スタックトレース取得という単一の目的に特化しており、非常にシンプルです。
  - **学習コストの低さ**: 新しい外部ライブラリのAPIや思想を学ぶ必要がありません。
- **劣っている点 (自前実装が)**:
  - **機能の豊富さ**: `cockroachdb/errors` は、スタックトレース以外にも、エラーレポート、エラーのグルーピング、ネットワーク経由でのエラー情報の伝達など、遥かに高機能です。
  - **コミュニティの知見**: 活発にメンテナンスされており、多くのエッジケースが考慮されています。自前実装では見逃している問題があるかもしれません。

### vs ログ出力時に都度スタックを取得する実装

- **優れている点 (自前実装が)**:
  - **正確なエラー発生地点**: エラーが **生成された時点** のスタックトレースを正確に記録できます。これはデバッグにおいて最も重要な情報です。
- **劣っている点 (自前実装が)**:
  - (特になし) ログ出力時にスタックを取得する方法では、エラー発生場所ではなく **ログが出力された場所** のスタックしか取れず、デバッグ情報としての価値が著しく低いため、我々の実装が明確に優れています。

## 5. 結論・総合評価

現在の自前実装は、**「`pkg/errors`への依存をなくし、Goの標準機能と`zerolog`の拡張性を最大限活用する」** という目的を見事に達成しています。

`errs.NewErr` の使用を規約とする必要はありますが、その見返りとして、**外部依存が少なく、軽量で、挙動が明確なエラーハンドリング基盤**を構築できています。`cockroachdb/errors` のような高機能ライブラリは、本プロジェクトの現在の要件に対してはオーバースペックであり、自前実装のシンプルさはメリットと言えます。

総じて、現在の実装はプロジェクトのニーズに合致した、**パフォーマンスとメンテナンス性のバランスが取れた、非常に実践的で優れたアプローチ**であると評価できます。
