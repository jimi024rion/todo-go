# DIツール選定 調査記録

このドキュメントは、DIツール選定のためのADR(Architecture Decision Record)作成に向けた、調査内容や比較検討の過程を記録するためのものです。

## 1. 評価軸

### 1.1. 技術的側面

- **解決タイミング**: コンパイル時に依存関係を検証・解決するか、実行時にリフレクションを用いて解決するか。
- **パフォーマンス**: 実行時オーバーヘッドの有無、アプリケーションの起動速度への影響。
- **エラー検出タイミング**: 依存関係の不整合などのエラーをコンパイル時に検出できるか、実行時まで分からないか。

### 1.2. 開発体験 (Developer Experience - DX)

- **学習コスト**: 新しい開発者がツールの概念や使い方を習得するのにかかる時間。
- **記述量 / ボイラープレート**: DI設定のためのコードや設定ファイルの記述がどれだけ必要か。
- **デバッグの容易さ**: 依存関係の解決に失敗した際や、意図しない挙動の際に原因を特定しやすいか。
- **"魔法"っぽさ / 透明性**: ツールの内部動作や依存関係の解決ロジックがどの程度明確に理解できるか。

### 1.3. プロジェクトとエコシステム

- **移行コスト**: 既存の`wire`ベースのコードから、新しいツールへ移行する際の作業量や難易度。
- **コミュニティ / 採用実績**: 困った時に情報（ドキュメント、サンプルコード、Q&Aなど）を見つけやすいか。ツールが広く使われているか。
- **メンテナンス状況**: プロジェクトが活発に開発・メンテナンスされているか（アーカイブ済みでないか、Issue/PRの対応状況など）。
- **機能セット**: 基本的なDI機能に加え、ライフサイクル管理、非同期初期化、リクエストスコープなどの高度な機能を提供するか。

## 2. 検討の背景

現在、プロジェクトではDIコンテナとして `google/wire` を利用している。
しかし、`google/wire` はGoogleによって **アーカイブ済み** となっており、今後のアップデートやメンテナンスが見込めない。
このため、将来的な技術的負債となるリスクを回避し、プロジェクトの持続可能性を高めるために代替ライブラリを検討する。

## 3. 要件

- **コンパイル時DI**: 実行時ではなくコンパイル時に依存関係のチェックを行い、型安全性を確保したい。リフレクションベースのライブラリも比較対象として調査するが、コンパイル時DIを優先する。
- **パフォーマンス**: 実行時のパフォーマンスへの影響が少ないこと。
- **学習コスト/移行コスト**: 既存の`wire`からの移行が比較的容易で、新しい開発者の学習コストが低いこと。
- **メンテナンス性**: コミュニティが活発で、継続的にメンテナンスされていること。

---

## 4. 候補ライブラリの調査

### 4.1. `google/wire` (現状維持)

- **分類**: コンパイル時DI
- **特徴**:
  - コード生成により、リフレクションなしで依存関係を解決。
  - コンパイル時に依存関係の欠落を検出可能。
  - 実行時パフォーマンスへの影響がほぼない。
- **懸念点**:
  - プロジェクトが[アーカイブ済み](https://github.com/google/wire)であり、今後のバグ修正や機能追加は期待できない。

### 4.2. `mazrean/kessoku` (有力候補)

- **分類**: コンパイル時DI
- **特徴**:
  - `wire`と同様にコード生成ベース。リフレクションを使用しないため、高いパフォーマンスと型安全性を両立。
  - **並列初期化**: 依存関係のないコンポーネントの初期化を並列実行できる `Async()` オプションがあり、アプリケーションの起動時間短縮に貢献する可能性がある。
- **`wire`との比較**:
  - APIの設計思想（`Provide`, `Inject`, `Bind`など）が`wire`に非常に似ており、学習・移行コストは低いと予想される。
  - `wire`にはない非同期初期化の機能が強み。
- **懸念点**:
  - `wire`や`fx`と比較してコミュニティが小さく、採用実績もまだ少ない。

### 4.3. `uber-go/dig` (比較対象)

- **分類**: 実行時DI (リフレクションベース)
- **特徴**:
  - Uberが開発したDIコンテナで、`uber-go/fx`の基盤となっている。
  - コンストラクタ関数をコンテナに登録し、`Invoke`で必要な依存関係を自動的に解決・注入する。
  - `go generate`のような事前コード生成は不要。
- **懸念点**:
  - **リフレクションを利用するため、依存関係の解決エラーは実行時まで検出できない。**
  - コンパイル時の型安全性に欠ける。
  - リフレクションによるわずかなパフォーマンスオーバーヘッドがある。
  - 依存関係が複雑になると、デバッグが困難になる場合がある。

### 4.4. `uber-go/fx` (比較対象)

- **分類**: 実行時DI (リフレクションベース) / アプリケーションフレームワーク
- **特徴**:
  - `dig`を内包し、さらにアプリケーションのライフサイクル管理（起動・シャットダウン）機能を提供する。
  - `fx.Lifecycle`フックにより、コンポーネントの初期化と後処理を宣言的に記述できる。
  - 大規模なアプリケーションでの実績が豊富。
- **懸念点**:
  - `dig`と同様、リフレクションベースであり、コンパイル時の型安全性は保証されない。
  - 学習コストが高く、アプリケーション全体が`fx`の流儀に強く依存する構造になる。
  - 導入すると、DIライブラリの選定というよりは「`fx`というフレームワークを使うか」という決定になる。

### 4.5. 手動DI (ライブラリ不使用)

- **分類**: DIパターン
- **特徴**:
  - ライブラリに依存せず、`main.go`などで依存関係を明示的に組み立てる。
  - `wire`のやっていることを手動で行うイメージ。
- **メリット**:
  - 最高の透明性とデバッグの容易さ。
  - 依存関係がゼロ。
  - 完全なコンパイル時安全性。
- **デメリット**:
  - アプリケーションが大きくなるにつれて、初期化コードが非常に長大で複雑になる。

---

## 5. 比較評価表

| 評価軸                        | `google/wire` (現状維持) | `mazrean/kessoku`   | `uber-go/dig` | `uber-go/fx`                | 手動DI       |
| :---------------------------- | :----------------------- | :------------------ | :------------ | :-------------------------- | :----------- |
| **解決タイミング**            | コンパイル時             | コンパイル時        | 実行時        | 実行時                      | コンパイル時 |
| **実行時オーバーヘッド**      | ほぼなし                 | ほぼなし            | 小さい        | 小さい                      | なし         |
| **起動パフォーマンス**        | 普通                     | **高速化の可能性**  | 普通          | 普通                        | 普通         |
| **エラー検出タイミング**      | コンパイル時             | コンパイル時        | 実行時        | 実行時                      | コンパイル時 |
| **学習コスト**                | 中                       | 中〜高              | 中            | 高                          | 低           |
| **記述量 / ボイラープレート** | 中 (コード生成)          | 中 (コード生成)     | 中            | 高                          | 高 (手動)    |
| **デバッグの容易さ**          | 高                       | 高                  | 中            | 中                          | 高           |
| **"魔法"っぽさ / 透明性**     | 中                       | 中                  | 中            | 高                          | 高           |
| **移行コスト**                | N/A                      | 低                  | 高            | 高                          | 中           |
| **コミュニティ / 採用実績**   | 大                       | 小                  | 大            | 大                          | N/A          |
| **メンテナンス状況**          | アーカイブ済み           | 活発                | 活発          | 活発                        | N/A          |
| **機能セット**                | 基本DI                   | 基本DI + 並列初期化 | 基本DI        | 基本DI + ライフサイクル管理 | 基本DI       |

---

## 6. 最終推奨

### 推奨: `mazrean/kessoku`

- **理由**:
  - ユーザーの最重要要件である**コンパイル時DI**と**実行時オーバーヘッドの少なさ**を完全に満たしている。
  - `google/wire`とAPIが酷似しており、学習コストと移行コストが非常に低いと見込まれる。
  - `wire`にはない**並列初期化**という明確な性能的メリット（特に起動パフォーマンスの改善）があり、これはアプリケーションの特性によっては大きな強みとなる。
  - `wire`がアーカイブされたことで生じる保守性の懸念を解消できる。

- **懸念点と対応案**:
  - **コミュニティの小ささ**: これが唯一の大きな懸念点となる。
    - **対応案1**: 導入前に、GitHubの活動状況（コミット頻度、IssueやPRへの対応速度）、ロードマップなどをより詳細に調査し、プロジェクトの安定性を見極める。
    - **対応案2**: `kessoku`のコードベースがシンプルで、自社で必要最低限のメンテナンスができる範囲であれば、コミュニティの大小よりもそのライブラリが提供する品質と機能（並列初期化）を優先する。
    - **対応案3**: 万が一の場合に備え、手動DIへの切り替えパスを念頭に置く。

### 次善の策/緊急時の選択肢: 手動DI

- **理由**: `kessoku`のコミュニティリスクが許容できないと判断された場合、手動DIは最も安全で透明性の高い選択肢である。DIライブラリに依存しないため、外部要因によるリスクは最小限になる。ただし、プロジェクト規模の拡大に伴う冗長性増加のリスクがある。

### 不採用とする候補

- **`google/wire`**: アーカイブ済みであるという代替検討の動機そのものであるため、不採用。
- **`uber-go/dig`**: ユーザーの「コンパイル時にエラー検出したい」という最重要要件を満たせないため、不採用。
- **`uber-go/fx`**: ユーザーの「コンパイル時にエラー検出したい」という最重要要件を満たせないため、不採用。また、DIライブラリの選定というよりは「フレームワークの採用」に近い決定となり、現状のプロジェクト構造との整合性も低い。

---

## 7. 簡易実装例 (mazrean/kessoku)

`mazrean/kessoku` を用いた簡単なDIの実装例を以下に示します。
この例では、`Logger` と `Config` の依存関係を `Service` に注入し、`Config` の初期化を非同期で行う様子を示します。

```go
// main.go (kessokuのサンプル)
package main

import (
 "fmt"
 "log"
 "time"

 "github.com/mazrean/kessoku"
)

// --- Define interfaces and structs ---
type Logger interface {
 Log(msg string)
}

type ConsoleLogger struct{}

func NewConsoleLogger() Logger {
 fmt.Println("ConsoleLogger initialized.")
 return &ConsoleLogger{}
}

func (l *ConsoleLogger) Log(msg string) {
 log.Println("[APP]", msg)
}

type Config struct {
 AppName string
 Version string
}

func NewConfig() *Config {
 // 非同期処理をシミュレートするため、少し待機
 time.Sleep(100 * time.Millisecond)
 fmt.Println("Config initialized.")
 return &Config{
  AppName: "MyApp",
  Version: "1.0.0",
 }
}

type Service struct {
 logger Logger
 config *Config
}

func NewService(logger Logger, config *Config) *Service {
 fmt.Println("Service initialized.")
 return &Service{
  logger: logger,
  config: config,
 }
}

// --- kessoku Providers ---
// kessoku.Provide() で通常のプロバイダを定義
var LoggerProvider = kessoku.Provide(NewConsoleLogger)

// kessoku.Async() で非同期実行可能なプロバイダを定義
// 依存関係がなければ、他のAsyncプロバイダと並列に実行される
var ConfigProvider = kessoku.Async(NewConfig)

// ServiceProvider は LoggerProvider と ConfigProvider に依存
var ServiceProvider = kessoku.Provide(NewService)

// --- kessoku Injector (go generate でコード生成される) ---
// Injectorを生成するための宣言
//go:generate kessoku inject
func InitApp() (*Service, error) {
 return kessoku.Inject[
  *Service, // 最終的に注入したい型
 ](
  // 利用するプロバイダを列挙
  LoggerProvider,
  ConfigProvider,
  ServiceProvider,
 )
}

func main() {
 fmt.Println("Starting application...")

 // 生成されたInitApp()を呼び出す
 // コンパイル時に依存関係が解決され、不正な場合はエラーになる
 appService, err := InitApp()
 if err != nil {
  log.Fatalf("Failed to initialize app: %v", err)
 }

 appService.logger.Log(fmt.Sprintf("%s (Version %s) started!", appService.config.AppName, appService.config.Version))
 fmt.Println("Application running.")
}
```

---

## 8. 簡易実装例 (google/wire)

`google/wire` を用いたDIの実装例です。`kessoku` と非常に似ていますが、`wire` の基本的な使い方を示します。

### 8.1. 準備

`wire.go` ファイルにプロバイダとインジェクタを定義します。

```go
// wire.go
//go:build wireinject
// +build wireinject

package main

import "github.com/google/wire"

// --- Define interfaces and structs (kessokuの例と同じ) ---
type Message string

func NewMessage() Message {
 return "Hello, World!"
}

type Greeter struct {
 Message Message
}

func NewGreeter(m Message) Greeter {
 return Greeter{Message: m}
}

func (g Greeter) Greet() Message {
 return g.Message
}

type Event struct {
 Greeter Greeter
}

func NewEvent(g Greeter) Event {
 return Event{Greeter: g}
}

func (e Event) Start() {
 msg := e.Greeter.Greet()
 println(msg)
}


// --- wire ProviderSet ---
// プロバイダをセットとしてまとめる
var SuperSet = wire.NewSet(NewMessage, NewGreeter, NewEvent)

// --- wire Injector ---
// インジェクタは生成したいオブジェクトを返す関数として定義
func InitializeEvent() (Event, error) {
 // wire.Build にプロバイダセットを渡す
 // 戻り値の error は未使用でも定義が必要
 wire.Build(SuperSet)
 return Event{}, nil // この部分はコード生成時に置き換えられる
}
```

### 8.2. コード生成

`wire` CLI を使ってコードを生成します。

```sh
go get github.com/google/wire/cmd/wire
wire
```

`wire_gen.go` というファイルが生成されます。

```go
// wire_gen.go
// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

// Injectors from wire.go:

func InitializeEvent() (Event, error) {
 message := NewMessage()
 greeter := NewGreeter(message)
 event := NewEvent(greeter)
 return event, nil
}
```

### 8.3. Goアプリケーションでの利用

```go
// main.go
package main

import "fmt"

// --- kessokuの例と同じインターフェースと構造体定義 ---
// ... (Message, Greeter, Event の定義は wire.go と同じ) ...

func main() {
 fmt.Println("Starting application with wire...")

 // 生成された InitializeEvent() を呼び出す
 event, err := InitializeEvent()
 if err != nil {
  panic(err)
 }

 event.Start()
 fmt.Println("Application running.")
}
```

---

## 9. 簡易実装例 (uber-go/dig)

`uber-go/dig` を用いたリフレクションベースDIの実装例です。

```go
package main

import (
 "fmt"
 "log"

 "go.uber.org/dig"
)

// --- Define interfaces and structs ---
type Config struct {
 AppName string
}

func NewConfig() (*Config, error) {
 fmt.Println("Initializing Config...")
 return &Config{AppName: "MyDigApp"}, nil
}

type Logger struct {
 AppName string
}

// LoggerはConfigに依存
func NewLogger(cfg *Config) *Logger {
 fmt.Println("Initializing Logger...")
 return &Logger{AppName: cfg.AppName}
}

func (l *Logger) Log(msg string) {
 log.Printf("[%s] %s", l.AppName, msg)
}

type Server struct {
 logger *Logger
}

// ServerはLoggerに依存
func NewServer(logger *Logger) *Server {
 fmt.Println("Initializing Server...")
 return &Server{logger: logger}
}

func (s *Server) Start() {
 s.logger.Log("Server started!")
}


func main() {
 // 1. コンテナを作成
 container := dig.New()

 // 2. プロバイダ (コンストラクタ) をコンテナに登録
 // Provideに失敗した場合はエラーハンドリング
 if err := container.Provide(NewConfig); err != nil {
  panic(err)
 }
 if err := container.Provide(NewLogger); err != nil {
  panic(err)
 }
 if err := container.Provide(NewServer); err != nil {
  panic(err)
 }

 // 3. 依存関係を解決して関数を実行
 // Invokeに渡した関数が必要とする引数が自動的に注入される
 err := container.Invoke(func(server *Server) {
  fmt.Println("Invoking function...")
  server.Start()
 })

 if err != nil {
  // 依存関係の不足や循環参照はここで実行時エラーになる
  panic(fmt.Sprintf("Failed to invoke server: %v", err))
 }

 fmt.Println("Application finished.")
}
```

---

## 10. 簡易実装例 (uber-go/fx)

`uber-go/fx` を用いたDIとアプリケーションライフサイクル管理の実装例です。

```go
package main

import (
 "context"
 "fmt"
 "log"
 "net/http"
 "time"

 "go.uber.org/fx"
)

// --- Logger ---
func NewLogger() *log.Logger {
 return log.New(log.Writer(), "[FX-APP] ", log.LstdFlags)
}

// --- HTTP Handler ---
func NewHandler(logger *log.Logger) (http.Handler, error) {
 logger.Println("Executing NewHandler.")
 mux := http.NewServeMux()
 mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  logger.Println("Got a request.")
  fmt.Fprintln(w, "Hello from fx!")
 })
 return mux, nil
}

// --- HTTP Server ---
func NewServer(lc fx.Lifecycle, logger *log.Logger, handler http.Handler) *http.Server {
 logger.Println("Executing NewServer.")
 server := &http.Server{
  Addr:    ":8080",
  Handler: handler,
 }

 // fxのライフサイクルにフックを登録
 lc.Append(fx.Hook{
  OnStart: func(ctx context.Context) error {
   logger.Println("Starting HTTP server.")
   // サーバーの起動はノンブロッキングで行う
   go func() {
    if err := server.ListenAndServe(); err != http.ErrServerClosed {
     logger.Fatalf("Failed to start server: %v", err)
    }
   }()
   return nil
  },
  OnStop: func(ctx context.Context) error {
   logger.Println("Stopping HTTP server.")
   // タイムアウト付きで優雅なシャットダウン
   shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
   defer cancel()
   return server.Shutdown(shutdownCtx)
  },
 })

 return server
}

func main() {
 // fx.Appを構築
 app := fx.New(
  // プロバイダを登録
  fx.Provide(
   NewLogger,
   NewHandler,
   NewServer,
  ),
  // アプリケーションのメインロジックを登録
  // ここではサーバーを起動するだけなので、引数で受け取るだけ
  fx.Invoke(func(*http.Server) {}),
 )

 // アプリケーションを実行
 // StartはCtrl+Cなどで停止されるまでブロックする
 app.Run()

 // アプリケーションが停止すると、OnStopフックが呼ばれる
 log.Println("Application stopped.")
}
```

---

## 11. 簡易実装例 (手動DI)

ライブラリを使わず、手動で依存関係を注入する例です。

```go
package main

import "fmt"

// --- 依存関係の定義 ---

type Database struct {
 ConnectionString string
}

func NewDatabase(connStr string) *Database {
 fmt.Printf("Connecting to database: %s\n", connStr)
 return &Database{ConnectionString: connStr}
}

type UserRepository struct {
 DB *Database
}

func NewUserRepository(db *Database) *UserRepository {
 fmt.Println("Creating user repository")
 return &UserRepository{DB: db}
}

type UserService struct {
 Repo *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
 fmt.Println("Creating user service")
 return &UserService{Repo: repo}
}


// --- アプリケーションの組み立て役 ---

// App は全てのコンポーネントを保持する
type App struct {
 UserService *UserService
}

// NewApp がDIコンテナの役割を果たす
func NewApp(connStr string) (*App, error) {
 // 依存関係のグラフを手動で構築する
 // 依存の末端から順に生成していく
 db := NewDatabase(connStr)
 userRepo := NewUserRepository(db)
 userService := NewUserService(userRepo)

 // 全てをApp構造体にまとめる
 return &App{
  UserService: userService,
 }, nil
}


func main() {
 fmt.Println("--- Manual Dependency Injection ---")

 // アプリケーションを初期化
 // 依存関係の解決は NewApp の中で行われる
 app, err := NewApp("postgres://user:pass@host/db")
 if err != nil {
  panic(err)
 }

 // アプリケーションのサービスを利用
 fmt.Println("App initialized, user service is ready.")
 _ = app.UserService // UserServiceが利用可能

 fmt.Println("--- Finished ---")
}
```
