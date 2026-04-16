# ── プロバイダー設定 ────────────────────────────────────

provider "google" {
  project = var.gcp_project_id
  region  = var.gcp_region
}

provider "neon" {
  api_key = var.neon_api_key
}

provider "vercel" {
  api_token = var.vercel_token
  # 個人アカウントの場合は team_id を省略できる
  team = var.vercel_team_id != "" ? var.vercel_team_id : null
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

  # Secret Manager に保存する DB 認証情報（Neon から取得）
  db_secrets = {
    DB_HOST     = neon_project.main.database_host
    DB_USER     = neon_role.app.name
    DB_PASSWORD = neon_role.app.password
    DB_NAME     = neon_database.main.name
  }

  # Atlas migration・GitHub Secrets 用の接続文字列
  database_url = "postgres://${neon_role.app.name}:${neon_role.app.password}@${neon_project.main.database_host}/${neon_database.main.name}?sslmode=require"
}
