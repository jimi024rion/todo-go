# ── プロバイダー設定 ────────────────────────────────────

provider "google" {
  project = var.gcp_project_id
  region  = var.gcp_region
}

provider "neon" {
  api_key = var.neon_api_key
}

# ── プロジェクト情報の取得 ───────────────────────────────
# project_number（Cloud Run のデフォルトSA算出に使用）を動的に取得する
data "google_project" "current" {
  project_id = var.gcp_project_id
}

# ── ローカル値 ───────────────────────────────────────────
locals {
  # Cloud Run が Secret Manager にアクセスするために使うデフォルト SA
  # 形式: {project_number}-compute@developer.gserviceaccount.com
  cloudrun_sa = "${data.google_project.current.number}-compute@developer.gserviceaccount.com"

  # Secret Manager に保存するのはパスワードのみ（他は平文で Cloud Run に直接セット）
  db_secrets = {
    DB_PASSWORD = neon_role.app.password
  }

  # Atlas migration・GitHub Secrets 用の接続文字列
  database_url = "postgres://${neon_role.app.name}:${neon_role.app.password}@${neon_project.main.database_host}/${neon_database.main.name}?sslmode=require"
}
