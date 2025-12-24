variable "db_url" {
  type    = string
  default = "postgres://user:P@ssw0rd@localhost:5432/todo-go?search_path=public&sslmode=disable"
}

env "local" {
  # ソース：SQLファイル（schema.sql）を状態の定義として使用
  src = "file://schema.sql"

  # 接続先：マイグレーションを適用するターゲットDB
  url = var.db_url

  # マイグレーションファイルの保存先
  migration {
    dir = "file://migrations"
  }

  # 差分計算用のテンポラリDB（Dev Database）
  # Atlasは正確な差分を出すために、一時的なクリーンなDBを必要とします
  dev = "docker://postgres/15/dev"
}
