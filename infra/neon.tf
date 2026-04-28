# ═══════════════════════════════════════════════════════
# Neon PostgreSQL（サーバーレスDB）
# ═══════════════════════════════════════════════════════
# Neon は PostgreSQL 互換のサーバーレス DB サービス
# 使用しない時間帯は自動的にスリープして無料枠に収まる

resource "neon_project" "main" {
  name      = "todo-go"
  org_id    = var.neon_org_id
  region_id = "aws-ap-southeast-1" # シンガポール（利用可能な最寄りリージョン）

  # 無料プランの上限値（21600秒 = 6時間）
  history_retention_seconds = 21600

  # 無料プランの最小スペックで起動
  default_endpoint_settings {
    autoscaling_limit_min_cu = 0.25 # 最小 0.25 vCPU
    autoscaling_limit_max_cu = 0.25 # 最大 0.25 vCPU（固定）
  }
}

# アプリ専用の DB ロール（ユーザー）
# パスワードは Neon が自動生成し、neon_role.app.password で参照できる
resource "neon_role" "app" {
  project_id = neon_project.main.id
  branch_id  = neon_project.main.default_branch_id
  name       = "app_user"
}

# アプリ用データベース
resource "neon_database" "main" {
  project_id = neon_project.main.id
  branch_id  = neon_project.main.default_branch_id
  name       = "todo"
  owner_name = neon_role.app.name
}
