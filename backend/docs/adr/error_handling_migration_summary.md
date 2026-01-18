## エラー処理の移行アプローチまとめ

`pkg/errors` からの移行にあたり、主に2つの有効なアプローチが考えられます。それぞれに長所・短所があるため、プロジェクトの方針に合わせて選択するのが理想的です。

### パターン1: `runtime`パッケージによる自前実装アプローチ

外部のエラー処理ライブラリに依存せず、Goの標準 `runtime` パッケージを利用して、プロジェクト固有のエラー型を拡張するアプローチです。

#### 概要

- **方針**: `errs.Err` 構造体にスタックトレース情報を直接保持させ、ロガー(`zerolog`)がそれを解釈して構造化ログを生成する。
- **責務**: エラー情報の保持からログ形式への変換まで、すべてプロジェクト内のコードで完結させる。

#### 主要な変更点

- **`errs.go`**:
  - `Err` 構造体に `callers []uintptr` フィールドを追加。
  - `NewErr`関数内で `runtime.Callers()` を呼び出し、スタックのプログラムカウンター情報を取得・保存する。
  - `StackFrames() []uintptr` のようなメソッドを実装し、保持しているスタック情報をロガーに公開する。
- **`logger.go`**:
  - `zerolog.ErrorStackMarshaler` を自前で実装する。
  - この中で `runtime.CallersFrames` を使い、プログラムカウンターからファイル名や行番号を解決し、`map`のスライス（JSON配列のもと）を生成する。
- **依存関係**:
  - 新たな外部ライブラリへの依存は発生しない。

#### メリット

- ✅ **依存の少なさ**: 外部依存がなく、プロジェクトが自己完結している。
- ✅ **透明性**: エラー処理のロジックがすべてプロジェクト内にあり、何が行われているか追跡しやすい。
- ✅ **構造の維持**: 元のコード構造からの変化が比較的小さい。

#### デメリット

- ❌ **メンテナンスコスト**: スタックトレースのフォーマットなど、エラー処理の基盤部分を自前で実装・メンテナンスし続ける必要がある。
- ❌ **拡張性**: 将来、より高度なエラー情報（任意のタグなど）を追加したくなった場合、`Err`構造体と関連処理を自力で拡張し続ける必要がある。

#### 主要コードスニペット

```go
// errs.go
type Err struct {
    code    ResultCode
    err     error
    callers []uintptr
}

func NewErr(code ResultCode, err error) error {
    callers := make([]uintptr, 64)
    n := runtime.Callers(2, callers)
    return &Err{code, err, callers[:n]}
}

func (e *Err) StackFrames() []uintptr {
    return e.callers
}
```

```go
// logger.go
zerolog.ErrorStackMarshaler = func(err error) interface{} {
    st, ok := err.(interface{ StackFrames() []uintptr })
    if !ok { return nil }

    frames := runtime.CallersFrames(st.StackFrames())
    var stack []map[string]interface{}
    for {
        frame, more := frames.Next()
        stack = append(stack, map[string]interface{}{
            "func":   frame.Function,
            "line":   frame.Line,
            "source": filepath.Base(frame.File),
        })
        if !more { break }
    }
    return stack
}
```

---

### パターン2: `cockroachdb/errors` 併用アプローチ （最終的な実装）

エラー処理のバックエンドとして高機能な `cockroachdb/errors` を採用しつつ、プロジェクト固有のエラー型でその詳細をカプセル化（隠蔽）するアプローチです。

#### 概要

- **方針**: エラー情報の保持とスタックトレース取得は `cockroachdb/errors` に任せ、`errs.Err` はそれをラップするインターフェースとして機能する。
- **責務**: `cockroachdb/errors` がエラー処理のエンジンとなり、ロガーは `pkg/errors` との互換インターフェースを通じてスタック情報を取得する。

#### 主要な変更点

- **`errs.go`**:
  - `Err` 構造体の内部で、`cockroachdb/errors` でラップされた `error` を保持する。
  - `NewErr` で `cockroachdb/errors.Wrapf` を呼び出し、スタックトレースをキャプチャさせる。
  - `StackTrace() pkgerrors.StackTrace` メソッドを実装。内部のエラーが持つ `StackTrace()` を呼び出すことで、`cockroachdb/errors` の情報を外部に公開する。
- **`logger.go`**:
  - `zerolog.ErrorStackMarshaler` に、`zerolog/pkgerrors` が提供する `MarshalStack` を設定する。
  - `MarshalStack` が `errs.Err` の `StackTrace()` メソッドを自動で認識し、構造化ログを生成してくれる。
- **依存関係**:
  - `github.com/cockroachdb/errors` への依存を追加する。

#### メリット

- ✅ **信頼性とメンテナンス性**: エラー処理の基盤を、実績があり活発にメンテナンスされているライブラリに任せられる。
- ✅ **高い拡張性**: `cockroachdb/errors` が持つ豊富な機能（任意の詳細情報付与、エラーレポート生成など）を将来的に活用できる。
- ✅ **関心の分離**: `errs.go` がエラー処理の詳細を隠蔽するため、ロガー等の呼び出し側は `cockroachdb/errors` の存在を意識せずに済む。コードの見通しが良くなる。

#### デメリット

- ❌ **依存の増加**: 新たに `cockroachdb/errors` という外部ライブラリへの依存が増える。
- ❌ **複雑さ**: `pkg/errors` とのインターフェース互換性を利用するなど、一見すると少しトリッキーに感じる部分があるかもしれない。

#### 主要コードスニペット

```go
// errs.go
import (
    "github.com/cockroachdb/errors"
    pkgerrors "github.com/pkg/errors"
)

type Err struct {
    code ResultCode
    err  error // cockroachdb/errors でラップされたエラー
}
func NewErr(code ResultCode, err error) error {
    return &Err{ code: code, err: errors.Wrapf(err, "") }
}
func (e *Err) StackTrace() pkgerrors.StackTrace {
    type stackTracer interface{ StackTrace() pkgerrors.StackTrace }
    var st stackTracer
    if errors.As(e.err, &st) { return st.StackTrace() }
    return nil
}
```

```go
// logger.go
import "github.com/rs/zerolog/pkgerrors"
func InitializeLogger() {
    // ...
    zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}
func (l *Logger) ErrorLog(err error) {
    var e *errs.Err
    event := l.logger.Error().Stack()
    if errors.As(err, &e) { /* ... ResultCodeを取得 ... */ }
    event.Err(err).Msg("")
}
```
