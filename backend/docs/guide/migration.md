# マイグレーション手順書

## はじめに

`atlas`を用いてマイグレーションする手順を以下に記す。

## 手順

### 1. ディレクトリ構成

```shell
.
├── atlas.hcl          # Atlasの設定ファイル
├── schema.sql         # 理想のテーブル定義（状態）を記述
└── migrations/        # 生成されたマイグレーションファイルが保存されるディレクトリ
```

### 2. atlasの設定

[atlas.hcl](./../../internal/infrastructure/rdb/atlas.hcl) にDBへの接続情報やマイグレーションのディレクトリ設定を記述する。

### 3. テーブル定義

[schema.sql](./schema.sql) にテーブル定義のSQLを記述する。

### 4. 差分の作成 (`migrate diff`)

マイグレーションファイル（SQL）を生成する。

```shell
atlas migrate diff ${name_of_migration} --env local
```

### 5. マイグレーション実行 (`migrate apply`)

生成されたマイグレーションファイルを実際のデータベースに適用する。

```shell
atlas migrate apply --env local
```
