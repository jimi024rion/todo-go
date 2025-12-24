# ORM/SQLジェネレータ選定 調査記録

このドキュメントは、ORM/SQLジェネレータ選定のためのADR(Architecture Decision Record)作成に向けた、調査内容や比較検討の過程を記録するためのものです。

## 1. 評価軸

### 1.1. 生産性 (Productivity)

- **コード生成**: スキーマ定義やSQLクエリから、モデルやCRUD操作などのGoコードをどれだけ自動生成してくれるか。
- **クエリの書きやすさ**: 型安全なクエリビルダーの使いやすさ、生のSQLに近い記述が可能か。
- **学習コスト**: 新しい開発者がツールの概念や使い方を習得するのにかかる時間。

### 1.2. パフォーマンス (Performance)

- **実行時オーバーヘッド**: 生成されたコードの実行速度、リフレクションの有無。
- **クエリの最適化**: 生成されるSQLの効率性や、複雑なクエリ（JOIN、集約など）を記述する際の柔軟性。

### 1.3. 安全性 (Safety)

- **型安全性**: コンパイル時にクエリやモデルの型チェックがどれだけ厳密に行われるか。
- **SQLインジェクション耐性**: 安全なクエリ構築をどれだけ強制・支援してくれるか。

### 1.4. エコシステムとメンテナンス (Ecosystem & Maintenance)

- **コミュニティ**: 活発さ、ドキュメントの充実度、情報（サンプルコード、Q&Aなど）の見つけやすさ。
- **プロジェクトの活発度**: 継続的にメンテナンスされているか（Issue/PRの対応状況、最終コミット日時など）。
- **マイグレーション機能**: スキーママイグレーションツールとの連携や、ツールが直接マイグレーション機能を提供するか。

### 1.5. 柔軟性 (Flexibility)

- **SQLの直接実行**: 生成されたコードやクエリビルダーでは対応できない複雑なクエリを、生のSQLで記述・実行できるか。
- **拡張性**: カスタム型やロジックの追加が容易か。

---

## 2. 検討の背景

現在、プロジェクトではORMとして `sqlboiler` を利用している。
`sqlboiler` はGoのコードを自動生成し、型安全なORMとして機能するが、代替検討の背景として以下の点が挙げられる。

- `sqlboiler/bob`への移行を検討しているが、他にもコード生成ベースの有力なORM/SQLジェネレータが存在するため、多角的に比較検討を行いたい。
- ORM選定においても、コンパイル時の型安全性やパフォーマンスを重視する。

## 3. 要件

- **コード生成ベース**: リフレクションではなく、スキーマやSQLからGoコードを生成するタイプを好む。
- **コンパイル時安全性**: コンパイル時に可能な限り多くのエラー（型不一致、クエリ構文など）を検出できること。
- **パフォーマンス**: 実行時のオーバーヘッドが少ない、または高度なクエリ最適化が可能であること。
- **開発体験**: 学習コスト、クエリの記述しやすさ、デバッグの容易さ。
- **メンテナンス性**: コミュニティが活発で、継続的にメンテナンスされていること。

---

## 4. 候補ライブラリの調査

### 4.1. `sqlboiler` (現状維持)

### 4.2. `sqlc` (比較対象)

- **概要**: SQLクエリからGoコードを自動生成するツール。厳密にはORMではなく、型安全なSQLクライアントライブラリ。
- **特徴**:
  - `.sql`ファイルにSQLクエリを直接記述し、`sqlc generate`コマンドで対応するGoの関数とGoの構造体を生成する。
  - 生成されたコードは、Goの`database/sql`パッケージを利用するシンプルなコードで構成される。
  - スキーマ（`schema.sql`）とクエリを元に、入力・出力の型をGoのコードで提供するため、高い型安全性を実現。
  - 開発者はSQLを直接記述できるため、クエリの最適化やデータベース固有の機能利用の自由度が高い。
- **メリット**:
  - **高い型安全性**: コンパイル時にクエリの構文チェックや型チェックが行われ、Goコードとの整合性も保証される。
  - **パフォーマンス**: 生成されるコードは`database/sql`の薄いラッパーであり、リフレクションもORMレイヤーもないため、非常に高速。
  - **SQLのフルコントロール**: SQLを直接記述するため、複雑なクエリや最適化されたクエリを自由に書ける。
  - **学習コスト**: SQLが書ければGoの経験が浅くても比較的容易に習得可能。
- **デメリット**:
  - **ORM機能の欠如**: リレーションシップの自動解決（Join、Preloadなど）や、モデルのライフサイクルフックのようなORMが提供する機能は持たない。
  - **記述量**: 各クエリに対応するSQLを全て記述する必要があるため、単純なCRUD操作でもやや記述量が多くなる傾向がある。
  - マイグレーション機能は提供しない（別途ツールが必要）。

### 4.3. `Bob` (有力候補)

- **概要**: `sqlboiler`の後継として位置づけられ、クエリビルダー、ORMコード生成、ファクトリ生成、SQLクエリコード生成（`sqlc`に類似）など、多機能を提供する。
- **特徴**:
  - データベースファーストなアプローチで、既存のスキーマからGoコードを生成する。
  - クエリビルダーとして `Squirrel` に似た流暢なAPIを提供し、型安全ではないが柔軟なクエリ構築が可能。
  - ORMとして、スキーマに基づいてモデルとCRUDメソッドを生成。リレーションシップ（Eager Loadingなど）もサポート。
  - `sqlc`のように、生SQLから型安全なGoコードを生成する機能も持つ。この機能は`sqlc`のいくつかの課題を解決しようとしている。
- **メリット**:
  - **多機能性**: クエリビルダー、ORM、ファクトリ、SQLジェネレータを一つのプロジェクトで提供。
  - **柔軟なクエリ構築**: クエリビルダーとORMの両方を提供するため、Goコード内でのクエリ構築の柔軟性が高い。
  - **型安全性**: 生成されるORMコードやSQLジェネレータによるコードは型安全。
  - **`sqlc`に対する利点 (SQLクエリ生成の側面)**:
    - 動的な`IN (?)`句やバルクインサートの扱いがより容易。
    - 生成されたクエリを「クエリモディファイア」として再利用・拡張できる。
  - **ORM機能**: `sqlc`が持たないリレーションシップの解決やフック機能を提供。
- **デメリット**:
  - **学習コスト**: 多機能ゆえに、学習すべき概念やAPIが多くなる傾向がある。
  - **記述量**: スキーマ定義やクエリ定義に加えて、コード生成のための設定が必要。
  - **`sqlc`の代替とはならない**: `Bob`のSQLクエリ生成機能は`sqlc`に似ているが、`sqlc`ほどの特化度やシンプルさはない。両者はSQLの利用方法に関して異なる哲学を持つ。

### 4.4. `sqlboiler` (現状維持)

- **概要**: データベーススキーマからGoのORMコードを自動生成するツール。
- **特徴**:
  - DBスキーマを元に、Goの構造体、CRUDメソッド、クエリビルダーなどを生成。
  - 生成されるコードはリフレクションを使用せず、高い型安全性とパフォーマンスを持つ。
  - リレーションシップやフック、バリデーションなど、多くのORM機能を提供。
- **メリット**:
  - **高い型安全性とパフォーマンス**: コンパイル時に多くのエラーを検出でき、実行時オーバーヘッドが低い。
  - **多機能なORM**: スキーマからフルセットのORM機能を自動生成。
  - **強力なクエリビルダー**: 型安全なクエリビルダーで複雑なクエリも記述可能。
- **デメリット**:
  - **開発の停滞**: 公式リポジトリがアーカイブされ、Bobへの移行が推奨されている。
  - **柔軟性**: 生成されたコードのカスタマイズ性が限定的。生のSQLを直接書く場合は別途対応が必要。

### 4.5. `ent` (比較対象)

- **概要**: Facebook (Meta) によって開発された、グラフベースのGo向けORM。スキーマをGoのコードで定義し、それに基づいて型安全なAPIを生成するコードファーストのアプローチを取る。
- **特徴**:
    - **コードファースト**: データベーススキーマをGoのコードとして定義する。マイグレーションも自動生成可能。
    - **グラフベースのAPI**: `user.QueryPosts().QueryComments().All(ctx)` のように、リレーションを直感的にたどるクエリが書ける。
    - **静的型付け**: コード生成により、クエリ、リレーション、マイグレーションに至るまで高い型安全性を実現する。
- **メリット**:
    - 非常に高い型安全性とコンパイル時チェック。
    - グラフ構造のデータを扱う際のクエリが非常に書きやすい。
    - スキーマ定義からマイグレーションまで一貫した開発体験。
- **デメリット**:
    - **コードファースト**: 既存のデータベースからモデルを生成するデータベースファーストのアプローチが基本の `sqlboiler` や `Bob` とは思想が異なる。本プロジェクトの「DBスキーマが正」という思想と合わない。
    - 学習コストが他のORMに比べてやや高い可能性がある。

### 4.6. `gorm` (比較対象)

- **概要**: Goで最も人気のあるORMの一つ。リフレクションベースで動作し、コード生成を必要としない。
- **特徴**:
    - **リフレクションベース**: Goの構造体タグ(`gorm:"..."`)を元に、実行時にSQLを組み立てる。
    - **機能豊富**: Preload/Joins、トランザクション、マイグレーション、フックなど、ORMに期待される機能を網羅している。
- **メリット**:
    - 非常に大きなコミュニティと豊富なドキュメント。
    - コード生成ステップが不要なため、迅速に開発を開始できる。
- **デメリット**:
    - **リフレクションベース**: 本プロジェクトの要件である「コード生成ベース」「コンパイル時安全性」を満たさない。パフォーマンスへの懸念や、エラーが実行時まで発覚しないリスクがある。
    - **型安全性の欠如**: `db.Where("name = ?", "jinzhu")` のように文字列でクエリ条件を記述する場面が多く、タイプミスがコンパイル時に検出できない。

---

## 5. 比較評価表

| 評価軸 | `sqlboiler` | `sqlc` | `Bob` | `ent` | `gorm` |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **アプローチ** | DB First | DB First | DB First | Code First | Code First |
| **コード生成** | ◎ (スキーマから) | ○ (SQLから) | ◎ (スキーマ/SQLから) | ◎ (Goコードから) | ✕ (リフレクション) |
| **実行時オーバーヘッド** | ほぼなし | ほぼなし | ほぼなし | ほぼなし | あり |
| **型安全性** | ◎ (コンパイル時) | ◎ (コンパイル時) | ◎ (コンパイル時) | ◎ (コンパイル時) | △ (実行時) |
| **クエリAPI** | クエリビルダー | SQL | クエリビルダー + SQL | グラフAPI | メソッドチェーン |
| **学習コスト** | 中 | 低 | 中〜高 | 高 | 低〜中 |
| **プロジェクト活発度**| 低 (アーカイブ) | 活発 | 活発 | 活発 | 活発 |
| **ORM機能** | ◎ | ✕ | ◎ | ◎ | ◎ |
| **マイグレーション** | ✕ | ✕ | ✕ | ◎ | △ |

---

## 6. 最終推奨

### 推奨: `Bob`

- **理由**:
    - プロジェクトの最重要要件である**コード生成ベース**かつ**コンパイル時安全性**を完全に満たしている。
    - `sqlboiler`の後継として位置づけられており、既存プロジェクトからの**移行パスが明確**。
    - クエリビルダー、ORM、`sqlc`ライクなSQLコード生成といった**多機能性**を提供し、プロジェクトの多様なニーズに対応できる。
    - DBスキーマを正とする**データベースファースト**のアプローチがプロジェクトの方針と合致する。

- **懸念点**:
    - 多機能であるため、学習コストは `sqlc` 等のよりシンプルなツールに比べて高くなる可能性がある。

### 不採用とする候補

- **`sqlboiler`**: 公式に開発が停止しアーカイブされているため、長期的には代替が必要。
- **`sqlc`**: ORM機能を持たず、`sqlboiler` からの移行先としては機能不足の懸念がある。
- **`gorm`**: リフレクションベースであり、プロジェクトの最重要要件である「コード生成ベースによるコンパイル時の安全性とパフォーマンス」を満たさないため不採用。
- **`ent`**: 非常に優れたORMだが、スキーマをGoコードで定義する「コードファースト」のアプローチを取る。これは、DBスキーマ定義を正とする「データベースファースト」を基本方針とする本プロジェクトとは思想が合わないため不採用。

---

## 7. 簡易実装例 (Bob)

`Bob`を用いたGoアプリケーションでの基本的なCRUD操作のサンプルコードを以下に示します。
この例では、`users`テーブルに対する`Insert`, `Find`, `Update`, `Delete`操作を`Bob`のORM機能とクエリビルダーを組み合わせて実行します。

```go
package main

import (
 "context"
 "database/sql"
 "fmt"
 "log"
 "time"

 _ "github.com/lib/pq" // PostgreSQL driver
 "github.com/stephenafamo/bob"
 "github.com/stephenafamo/bob/dialect/psql"
 "github.com/stephenafamo/bob/dialect/psql/sm" // for query builder select mods
 "github.com/stephenafamo/bob/types"          // for nullable fields

 // bobgenによって生成されるモデルをインポートする前提
 // 実際には `your_project/backend/internal/infrastructure/rdb/models` など
 // このサンプルでは簡略化のため、`models`パッケージが存在すると仮定
 // また、`models.User` や `models.Users` は bobgen によって生成されるものとする
 // (実際の生成コードは複雑なため、ここでは概念的な表現に留める)
 // ---
 // package models
 //
 // import (
 //  "context"
 //  "github.com/stephenafamo/bob"
 //  "github.com/stephenafamo/bob/dialect/psql"
 //  "github.com/stephenafamo/bob/dialect/psql/um" // for update mods
 //  "github.com/stephenafamo/bob/dialect/psql/dm" // for delete mods
 //  "github.com/stephenafamo/bob/types"
 // )
 //
 // type User struct {
 //  ID        int64           `bob:"id,pk"`
 //  Name      string          `bob:"name"`
 //  Email     string          `bob:"email,unique"`
 //  CreatedAt types.Timestamp `bob:"created_at"`
 //  UpdatedAt types.Timestamp `bob:"updated_at"`
 // }
 //
 // // Users table struct (generated)
 // var Users = struct {
 //  Table psql.Table
 //  Columns struct {
 //   ID        psql.Column
 //   Name      psql.Column
 //   Email     psql.Column
 //   CreatedAt psql.Column
 //   UpdatedAt psql.Column
 //  }
 // }{
 //  Table: psql.NewTable("users"),
 //  Columns: struct {
 //   ID        psql.Column
 //   Name      psql.Column
 //   Email     psql.Column
 //   CreatedAt psql.Column
 //   UpdatedAt psql.Column
 //  }{
 //   ID:        psql.NewColumn("id"),
 //   Name:      psql.NewColumn("name"),
 //   Email:     psql.NewColumn("email"),
 //   CreatedAt: psql.NewColumn("created_at"),
 //   UpdatedAt: psql.NewColumn("updated_at"),
 //  },
 // }
 //
 // // UserモデルのCRUD操作メソッド群 (bobgenが生成する想定)
 // func (u *User) Insert(ctx context.Context, exec bob.Executor) error {
 //  // ... 実際の生成コード ...
 //  // ここでは簡略化し、psql.Insertを使用
 //  return psql.Insert(u).Into(Users.Table).Returning(psql.Stars).One(ctx, exec, u)
 // }
 // func (u *User) Update(ctx context.Context, exec bob.Executor) (int64, error) {
 //  // ... 実際の生成コード ...
 //  return psql.Update(u).Table(Users.Table).Set(
 //   Users.Columns.Name.Set(psql.Arg(u.Name)),
 //   Users.Columns.Email.Set(psql.Arg(u.Email)),
 //   Users.Columns.UpdatedAt.Set(psql.Arg(u.UpdatedAt)),
 //  ).Where(Users.Columns.ID.EQ(psql.Arg(u.ID))).Exec(ctx, exec)
 // }
 // func (u *User) Delete(ctx context.Context, exec bob.Executor) (int64, error) {
 //  // ... 実際の生成コード ...
 //  return psql.Delete().From(Users.Table).Where(Users.Columns.ID.EQ(psql.Arg(u.ID))).Exec(ctx, exec)
 // }
 //
 // func FindUser(ctx context.Context, exec bob.Executor, id int64) (*User, error) {
 //  var user User
 //  err := psql.Select(Users.Columns.ID, Users.Columns.Name, Users.Columns.Email, Users.Columns.CreatedAt, Users.Columns.UpdatedAt).
 //   From(Users).Where(Users.Columns.ID.EQ(psql.Arg(id))).One(ctx, exec, &user)
 //  if err == sql.ErrNoRows {
 //   return nil, nil // Or a custom not found error
 //  }
 //  return &user, err
 // }
 //
 // type userQuery struct{}
 // func (q userQuery) All(ctx context.Context, exec bob.Executor, mods ...bob.Mod[sm.SelectMod]) ([]*User, error) {
 //  var users []*User
 //  err := psql.Select(Users.Columns.ID, Users.Columns.Name, Users.Columns.Email, Users.Columns.CreatedAt, Users.Columns.UpdatedAt).
 //   From(Users).Apply(mods...).All(ctx, exec, &users)
 //  return users, err
 // }
 //
 // var UsersQuery userQuery
 // ---
)

// main関数のための簡易的なモデルとテーブル表現 (通常はbobgenが生成)
type User struct {
 ID        int64
 Name      string
 Email     string
 CreatedAt time.Time
 UpdatedAt time.Time
}

// 実際のbobgenで生成されるmodelsパッケージが存在する前提で記述
// このサンプルをそのまま実行するためには、modelsパッケージを別途定義するか
// bobgenで生成する必要があります。
// 以下の関数は、bobgenが生成する`models`パッケージの関数を模倣しています。
// 実際のBobの利用例では、これらの関数が自動生成されていると仮定してください。
// `bobgen.yaml`の設定に応じて、生成されるパッケージ名や構造体名、メソッド名が異なります。
// 例えば、`backend/internal/infrastructure/rdb/models` に生成されると仮定します。

// FindUser - bobgenが生成するFindUser関数を模倣
func FindUser(ctx context.Context, exec bob.Executor, id int64) (*User, error) {
 // 実際にはbobgenが生成するクエリビルダーが使われる
 // 便宜上、ここでは手動でSelect文を構築
 row := exec.QueryRowContext(ctx, "SELECT id, name, email, created_at, updated_at FROM users WHERE id = $1", id)
 var user User
 err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
 if err == sql.ErrNoRows {
  return nil, nil // Not found
 }
 return &user, err
}

// InsertUser - bobgenが生成するUser.Insert()メソッドを模倣
func InsertUser(ctx context.Context, exec bob.Executor, user *User) error {
 // 実際にはbobgenが生成するInsertメソッドが使われる
 row := exec.QueryRowContext(ctx, "INSERT INTO users (name, email, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id",
  user.Name, user.Email, user.CreatedAt, user.UpdatedAt)
 return row.Scan(&user.ID)
}

// UpdateUser - bobgenが生成するUser.Update()メソッドを模倣
func UpdateUser(ctx context.Context, exec bob.Executor, user *User) (int64, error) {
 // 実際にはbobgenが生成するUpdateメソッドが使われる
 result, err := exec.ExecContext(ctx, "UPDATE users SET name = $1, email = $2, updated_at = $3 WHERE id = $4",
  user.Name, user.Email, user.UpdatedAt, user.ID)
 if err != nil {
  return 0, err
 }
 return result.RowsAffected()
}

// DeleteUser - bobgenが生成するUser.Delete()メソッドを模倣
func DeleteUser(ctx context.Context, exec bob.Executor, id int64) (int64, error) {
 // 実際にはbobgenが生成するDeleteメソッドが使われる
 result, err := exec.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
 if err != nil {
  return 0, err
 }
 return result.RowsAffected()
}

// AllUsers - bobgenが生成するUsers.Query().All()を模倣
func AllUsers(ctx context.Context, exec bob.Executor) ([]*User, error) {
 rows, err := exec.QueryContext(ctx, "SELECT id, name, email, created_at, updated_at FROM users ORDER BY id")
 if err != nil {
  return nil, err
 }
 defer rows.Close()

 var users []*User
 for rows.Next() {
  var user User
  if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
   return nil, err
  }
  users = append(users, &user)
 }
 return users, rows.Err()
}


func main() {
 connStr := "user=postgres password=postgres dbname=bob_example host=localhost sslmode=disable"
 sqlDB, err := sql.Open("postgres", connStr)
 if err != nil {
  log.Fatalf("Error opening database: %v", err)
 }
 defer sqlDB.Close()

 // Wrap *sql.DB with bob.DB for bob's context and dialect features
 // Note: For bobgen generated models, you often pass bob.Executor directly.
 // This bDB object will be the primary executor.
 bDB := bob.NewDB(sqlDB)
 ctx := context.Background()

 // ---
 // Schema setup (simplified, in real app use migrations like Atlas)
 // ---
 _, err = bDB.Exec(ctx, `
  DROP TABLE IF EXISTS users CASCADE;
  CREATE TABLE users (
   id SERIAL PRIMARY KEY,
   name TEXT NOT NULL,
   email TEXT UNIQUE NOT NULL,
   created_at TIMESTAMPTZ DEFAULT NOW(),
   updated_at TIMESTMPTZ DEFAULT NOW()
  );
 `)
 if err != nil {
  log.Fatalf("Failed to setup schema: %v", err)
 }
 fmt.Println("Schema setup complete.")

 // ---
 // 1. Create (Insert)
 // ---
 fmt.Println("\n--- Create User ---")
 user1 := &User{
  Name:      "Alice",
  Email:     "alice@example.com",
  CreatedAt: time.Now(),
  UpdatedAt: time.Now(),
 }
 if err := InsertUser(ctx, bDB, user1); err != nil {
  log.Fatalf("Failed to insert user1: %v", err)
 }
 fmt.Printf("Inserted user: ID=%d, Name=%s\n", user1.ID, user1.Name)

 user2 := &User{
  Name:      "Bob",
  Email:     "bob@example.com",
  CreatedAt: time.Now(),
  UpdatedAt: time.Now(),
 }
 if err := InsertUser(ctx, bDB, user2); err != nil {
  log.Fatalf("Failed to insert user2: %v", err)
 }
 fmt.Printf("Inserted user: ID=%d, Name=%s\n", user2.ID, user2.Name)


 // ---
 // 2. Read (Find by ID)
 // ---
 fmt.Println("\n--- Find User by ID ---")
 foundUser, err := FindUser(ctx, bDB, user1.ID)
 if err != nil {
  log.Fatalf("Failed to find user by ID %d: %v", user1.ID, err)
 }
 fmt.Printf("Found user: %+v\n", foundUser)

 // ---
 // 3. Read (All users)
 // ---
 fmt.Println("\n--- List All Users ---")
 allUsers, err := AllUsers(ctx, bDB)
 if err != nil {
  log.Fatalf("Failed to list all users: %v", err)
 }
 for _, u := range allUsers {
  fmt.Printf("- ID: %d, Name: %s, Email: %s\n", u.ID, u.Name, u.Email)
 }

 // ---
 // 4. Update
 // ---
 fmt.Println("\n--- Update User ---")
 user1.Name = "Alice Smith"
 user1.UpdatedAt = time.Now()
 if _, err := UpdateUser(ctx, bDB, user1); err != nil {
  log.Fatalf("Failed to update user: %v", err)
 }
 fmt.Printf("Updated user: ID=%d, New Name=%s\n", user1.ID, user1.Name)

 // Verify update
 updatedUser, err := FindUser(ctx, bDB, user1.ID)
 if err != nil {
  log.Fatalf("Failed to verify updated user: %v", err)
 }
 fmt.Printf("Verified updated user: %+v\n", updatedUser)

 // ---
 // 5. Delete
 // ---
 fmt.Println("\n--- Delete User ---")
 if _, err := DeleteUser(ctx, bDB, user1.ID); err != nil {
  log.Fatalf("Failed to delete user: %v", err)
 }
 fmt.Printf("Deleted user: ID=%d\n", user1.ID)

 // Verify delete
 _, err = FindUser(ctx, bDB, user1.ID)
 if err == nil {
  log.Fatal("User was not deleted!")
 }
 fmt.Println("User successfully deleted (not found).")
}
```

---

## 8. 簡易実装例 (sqlc)

`sqlc` を用いたGoアプリケーションでの基本的なCRUD操作のサンプルを以下に示します。
`sqlc`は、開発者が書いたSQLクエリから型安全なGoのコードを生成します。

### 8.1. 準備

以下の3つのファイルを用意します。

1. **`sqlc.yaml`**: `sqlc` の設定ファイル。
2. **`schema.sql`**: データベースのテーブル定義。
3. **`query.sql`**: 実行したいSQLクエリ。

```yaml
# sqlc.yaml
version: "2"
sql:
  - engine: "postgresql"
    schema: "schema.sql"
    queries: "query.sql"
    gen:
      go:
        package: "authors"
        out: "authors"
```

```sql
-- schema.sql
CREATE TABLE authors (
  id   BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  bio  TEXT
);
```

```sql
-- query.sql
-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = $1 LIMIT 1;

-- name: ListAuthors :many
SELECT * FROM authors
ORDER BY name;

-- name: CreateAuthor :one
INSERT INTO authors (
  name, bio
) VALUES (
  $1, $2
)
RETURNING *;

-- name: UpdateAuthor :one
UPDATE authors
SET name = $2, bio = $3
WHERE id = $1
RETURNING *;

-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = $1;
```

### 8.2. コード生成

上記ファイルを用意した後、`sqlc generate` コマンドを実行すると、`authors` ディレクトリにGoのコードが生成されます。

```sh
sqlc generate
```

生成されるファイル:

- `authors/db.go`
- `authors/models.go`
- `authors/query.sql.go`

### 8.3. Goアプリケーションでの利用

生成されたコードをGoアプリケーションから利用する例です。

```go
package main

import (
 "context"
 "database/sql"
 "fmt"
 "log"

 _ "github.com/lib/pq"

 // sqlcが生成したパッケージをインポート
 "your_project/authors" // パスは適切に変更してください
)

func main() {
 connStr := "user=postgres password=postgres dbname=sqlc_example host=localhost sslmode=disable"
 db, err := sql.Open("postgres", connStr)
 if err != nil {
  log.Fatal(err)
 }
 defer db.Close()

 // スキーマのセットアップ
 _, err = db.Exec(`
  DROP TABLE IF EXISTS authors;
  CREATE TABLE authors (
   id   BIGSERIAL PRIMARY KEY,
   name TEXT NOT NULL,
   bio  TEXT
  );
 `)
 if err != nil {
  log.Fatalf("Failed to setup schema: %v", err)
 }

 ctx := context.Background()
 queries := authors.New(db) // 生成されたNew関数でQueriesインスタンスを作成

 // ---
 // 1. Create
 // ---
 fmt.Println("--- Create Author ---")
 insertedAuthor, err := queries.CreateAuthor(ctx, authors.CreateAuthorParams{
  Name: "Brian Kernighan",
  Bio:  sql.NullString{String: "Co-author of The C Programming Language", Valid: true},
 })
 if err != nil {
  log.Fatal(err)
 }
 fmt.Printf("Inserted Author: %+v\n", insertedAuthor)

 // ---
 // 2. Read (One)
 // ---
 fmt.Println("\n--- Get Author ---")
 author, err := queries.GetAuthor(ctx, insertedAuthor.ID)
 if err != nil {
  log.Fatal(err)
 }
 fmt.Printf("Found Author: %+v\n", author)

 // ---
 // 3. Update
 // ---
 fmt.Println("\n--- Update Author ---")
 updatedAuthor, err := queries.UpdateAuthor(ctx, authors.UpdateAuthorParams{
  ID:   author.ID,
  Name: "B. W. Kernighan",
  Bio:  author.Bio,
 })
 if err != nil {
  log.Fatal(err)
 }
 fmt.Printf("Updated Author: %+v\n", updatedAuthor)


 // ---
 // 4. Read (Many)
 // ---
 fmt.Println("\n--- List Authors ---")
 authorList, err := queries.ListAuthors(ctx)
 if err != nil {
  log.Fatal(err)
 }
 for _, a := range authorList {
  fmt.Printf("- Author: %+v\n", a)
 }

 // ---
 // 5. Delete
 // ---
 fmt.Println("\n--- Delete Author ---")
 err = queries.DeleteAuthor(ctx, author.ID)
 if err != nil {
  log.Fatal(err)
 }
 fmt.Printf("Deleted author with ID: %d\n", author.ID)

 // Verify delete
 authorList, err = queries.ListAuthors(ctx)
 if err != nil {
  log.Fatal(err)
 }
 fmt.Printf("Authors after deletion: %d\n", len(authorList))
}

```

---

## 9. 簡易実装例 (sqlboiler)

`sqlboiler` を用いたGoアプリケーションでの基本的なCRUD操作のサンプルを以下に示します。
`sqlboiler` はデータベーススキーマから型安全なGoのORMコードを生成します。

### 9.1. 準備

1. **データベースとテーブルの準備**: `sqlboiler`は既存のデータベーススキーマを読み取ってコードを生成します。

    ```sql
    -- 例: psql でデータベースに接続してテーブルを作成
    CREATE TABLE pilots (
        id   SERIAL PRIMARY KEY,
        name TEXT NOT NULL
    );

    CREATE TABLE jets (
        id          SERIAL PRIMARY KEY,
        pilot_id    INTEGER NOT NULL REFERENCES pilots(id),
        age         INTEGER NOT NULL,
        name        TEXT NOT NULL,
        color       TEXT NOT NULL
    );
    ```

2. **設定ファイル `sqlboiler.toml`**:

    ```toml
    # sqlboiler.toml
    [psql]
      dbname = "sqlboiler_example"
      host   = "localhost"
      port   = 5432
      user   = "postgres"
      pass   = "postgres"
      sslmode= "disable"
      output = "models"
      pkgname= "models"
    ```

### 9.2. コード生成

`sqlboiler psql` コマンドを実行すると、`models` ディレクトリにORMコードが生成されます。

```sh
sqlboiler psql
```

生成されるファイル:

- `models/pilots.go`
- `models/jets.go`
- `models/boil_types.go`
- `models/boil_table_names.go`
- etc...

### 9.3. Goアプリケーションでの利用

生成されたORMモデルをGoアプリケーションから利用する例です。

```go
package main

import (
 "context"
 "database/sql"
 "fmt"
 "log"

 "github.com/volatiletech/sqlboiler/v4/boil"
 "github.com/volatiletech/sqlboiler/v4/queries/qm"

 // sqlboilerが生成したパッケージをインポート
 "your_project/models" // パスは適切に変更してください

 _ "github.com/lib/pq"
)

func main() {
 connStr := "user=postgres password=postgres dbname=sqlboiler_example host=localhost sslmode=disable"
 db, err := sql.Open("postgres", connStr)
 if err != nil {
  log.Fatal(err)
 }
 defer db.Close()

 // グローバルなDBコネクションを設定 (sqlboilerの推奨)
 boil.SetDB(db)

 ctx := context.Background()

 // ---
 // 0. Clean up
 // ---
 // 関連を考慮して削除
 models.Jets().DeleteAll(ctx, db)
 models.Pilots().DeleteAll(ctx, db)


 // ---
 // 1. Create (Insert)
 // ---
 fmt.Println("--- Create Pilot & Jet ---")
 pilot := models.Pilot{Name: "John"}
 err = pilot.Insert(ctx, db, boil.Infer())
 if err != nil {
  log.Fatal(err)
 }
 fmt.Printf("Inserted Pilot: ID=%d, Name=%s\n", pilot.ID, pilot.Name)


 jet := models.Jet{
  PilotID: pilot.ID,
  Age:     20,
  Name:    "F-16",
  Color:   "Gray",
 }
 err = jet.Insert(ctx, db, boil.Infer())
 if err != nil {
  log.Fatal(err)
 }
 fmt.Printf("Inserted Jet: ID=%d, Name=%s, PilotID=%d\n", jet.ID, jet.Name, jet.PilotID)


 // ---
 // 2. Read (Eager Loading)
 // ---
 fmt.Println("\n--- Find Jet with Pilot ---")
 // qm.Loadを使用して関連するPilotをEager Loadする
 foundJet, err := models.Jets(
  models.JetWhere.ID.EQ(jet.ID),
  qm.Load(models.JetRels.Pilot),
 ).One(ctx, db)
 if err != nil {
  log.Fatal(err)
 }
 // foundJet.R.Pilot で関連先のデータにアクセスできる
 fmt.Printf("Found Jet: %s, Pilot: %s\n", foundJet.Name, foundJet.R.Pilot.Name)


 // ---
 // 3. Update
 // ---
 fmt.Println("\n--- Update Pilot ---")
 pilot.Name = "Captain John"
 rowsAff, err := pilot.Update(ctx, db, boil.Infer())
 if err != nil {
  log.Fatal(err)
 }
 if rowsAff != 1 {
  log.Fatalf("expected 1 row to be affected, got %d", rowsAff)
 }
 fmt.Printf("Updated Pilot: ID=%d, New Name=%s\n", pilot.ID, pilot.Name)

 // ---
 // 4. Read (QueryMods)
 // ---
 fmt.Println("\n--- List all Pilots ---")
 allPilots, err := models.Pilots().All(ctx, db)
 if err != nil {
  log.Fatal(err)
 }
 for _, p := range allPilots {
  fmt.Printf("- Pilot: ID=%d, Name=%s\n", p.ID, p.Name)
 }

 // ---
 // 5. Delete
 // ---
 fmt.Println("\n--- Delete Jet & Pilot ---")
 // 関連するJetを先に削除
 rowsAff, err = foundJet.Delete(ctx, db)
 if err != nil {
  log.Fatal(err)
 }
 fmt.Printf("Deleted Jet, rows affected: %d\n", rowsAff)

 rowsAff, err = pilot.Delete(ctx, db)
 if err != nil {
  log.Fatal(err)
 }
 fmt.Printf("Deleted Pilot, rows affected: %d\n", rowsAff)

 // Verify
 count, err := models.Pilots().Count(ctx, db)
 if err != nil {
  log.Fatal(err)
 }
 fmt.Printf("Pilots count after deletion: %d\n", count)
}
```
---

## 10. 簡易実装例 (ent)

`ent` を用いたGoアプリケーションでの基本的なCRUD操作のサンプルを以下に示します。
`ent` はGoコードでスキーマを定義し、コード生成を経て型安全な操作を実現します。

### 10.1. 準備

1.  **スキーマ定義**: `ent/schema` ディレクトリに、エンティティのスキーマをGoコードで記述します。

    ```go
    // ent/schema/user.go
    package schema

    import (
        "entgo.io/ent"
        "entgo.io/ent/schema/field"
    )

    type User struct {
        ent.Schema
    }

    func (User) Fields() []ent.Field {
        return []ent.Field{
            field.String("name").
                Default("unknown"),
            field.Int("age").
                Positive(),
        }
    }
    ```

### 10.2. コード生成

プロジェクトルートで `go run entgo.io/ent/cmd/ent init User` などを実行して初期化した後、`go generate ./ent` コマンドでコードを生成します。

```sh
$ go generate ./ent
```

これにより、`ent` ディレクトリにクライアントやモデルのコードが生成されます。

### 10.3. Goアプリケーションでの利用

生成されたクライアントを使ってアプリケーションを記述します。

```go
package main

import (
	"context"
	"fmt"
	"log"

	"your_project/ent" // 生成されたentパッケージをインポート

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// インメモリのSQLiteに接続するクライアントを作成
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	// スキーママイグレーションを実行
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	// ---
	// 1. Create
	// ---
	fmt.Println("--- Create User ---")
	u, err := client.User.
		Create().
		SetAge(30).
		SetName("a8m").
		Save(ctx)
	if err != nil {
		log.Fatalf("failed creating user: %v", err)
	}
	fmt.Printf("Created User: %+v\n", u)

	// ---
	// 2. Read
	// ---
	fmt.Println("\n--- Get User ---")
	foundUser, err := client.User.Get(ctx, u.ID)
	if err != nil {
		log.Fatalf("failed getting user: %v", err)
	}
	fmt.Printf("Found User: %+v\n", foundUser)

	// ---
	// 3. Update
	// ---
	fmt.Println("\n--- Update User ---")
	updatedUser, err := foundUser.Update().
		SetAge(31).
		Save(ctx)
	if err != nil {
		log.Fatalf("failed updating user: %v", err)
	}
	fmt.Printf("Updated User: %+v\n", updatedUser)

	// ---
	// 4. Query
	// ---
	fmt.Println("\n--- Query Users ---")
	users, err := client.User.Query().
		Where(ent.User.AgeGT(30)).
		All(ctx)
	if err != nil {
		log.Fatalf("failed querying users: %v", err)
	}
	fmt.Printf("Found Users (age > 30): %+v\n", users)

	// ---
	// 5. Delete
	// ---
	fmt.Println("\n--- Delete User ---")
	err = client.User.DeleteOneID(u.ID).Exec(ctx)
	if err != nil {
		log.Fatalf("failed deleting user: %v", err)
	}
	fmt.Printf("Deleted User with ID: %d\n", u.ID)

	// Verify
	count, err := client.User.Query().Count(ctx)
	if err != nil {
		log.Fatalf("failed counting users: %v", err)
	}
	fmt.Printf("User count after deletion: %d\n", count)
}
```

---

## 11. 簡易実装例 (gorm)

`gorm` を用いたリフレクションベースのORM操作のサンプルです。

### 11.1. モデル定義

Goの構造体とgormタグでモデルを定義します。

```go
// models.go
package main

import "gorm.io/gorm"

type Product struct {
	gorm.Model // ID, CreatedAt, UpdatedAt, DeletedAt を含む
	Code  string
	Price uint
}
```

### 11.2. Goアプリケーションでの利用

```go
package main

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// インメモリのSQLiteに接続
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	// スキーママイグレーションを実行
	db.AutoMigrate(&Product{})

	// ---
	// 1. Create
	// ---
	fmt.Println("--- Create Product ---")
	product := &Product{Code: "D42", Price: 100}
	db.Create(product)
	fmt.Printf("Created Product: ID=%d, Code=%s\n", product.ID, product.Code)


	// ---
	// 2. Read
	// ---
	fmt.Println("\n--- Get Product ---")
	var foundProduct Product
	// 主キーで検索
	db.First(&foundProduct, product.ID)
	fmt.Printf("Found Product: %+v\n", foundProduct)

	// 条件で検索
	var foundByCode Product
	db.First(&foundByCode, "code = ?", "D42")
	fmt.Printf("Found By Code: %+v\n", foundByCode)

	// ---
	// 3. Update
	// ---
	fmt.Println("\n--- Update Product ---")
	db.Model(&foundProduct).Update("Price", 200)
	fmt.Printf("Updated Product Price: %d\n", foundProduct.Price)

	// 複数のフィールドを更新
	db.Model(&foundProduct).Updates(Product{Price: 250, Code: "F42"})
	fmt.Printf("Updated Product: %+v\n", foundProduct)

	// ---
	// 4. Delete
	// ---
	fmt.Println("\n--- Delete Product ---")
	db.Delete(&foundProduct, foundProduct.ID)
	fmt.Printf("Deleted product with ID: %d\n", foundProduct.ID)

	// Verify
	var result Product
	err = db.First(&result, foundProduct.ID).Error
	if err != nil {
		fmt.Printf("Product not found after deletion: %v\n", err)
	}
}
```
