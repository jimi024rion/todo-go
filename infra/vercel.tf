# ═══════════════════════════════════════════════════════
# Vercel（フロントエンド ホスティング）
# ═══════════════════════════════════════════════════════

resource "vercel_project" "frontend" {
  name      = "todo-go-frontend"
  framework = "nextjs"

  # モノレポ構成のため frontend/ サブディレクトリを指定
  root_directory = "frontend"

  # main ブランチへの push で本番デプロイ
  git_repository = {
    type = "github"
    repo = var.github_repo
  }
}

# ── 環境変数 ──────────────────────────────────────────
# BACKEND_URL は Next.js の API Route（サーバー側）が参照する
# NEXT_PUBLIC_ プレフィックスがないためブラウザのバンドルに含まれない

resource "vercel_project_environment_variable" "backend_url_production" {
  project_id = vercel_project.frontend.id
  key        = "BACKEND_URL"
  value      = "https://${google_cloud_run_v2_service.backend.uri}"
  target     = ["production"]
}

resource "vercel_project_environment_variable" "backend_url_preview" {
  project_id = vercel_project.frontend.id
  key        = "BACKEND_URL"
  # プレビュー環境（PR ごとのデプロイ）はローカルバックエンドを参照
  # 実際の開発では Preview 用の Cloud Run サービスを用意するのが理想
  value  = "http://localhost:8080"
  target = ["preview"]
}
