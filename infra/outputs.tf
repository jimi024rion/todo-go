# ═══════════════════════════════════════════════════════
# GitHub Actions に設定する値
# ═══════════════════════════════════════════════════════
# terraform apply 後に以下のコマンドで確認できる:
#   terraform output github_variables    # 非機密値
#   terraform output -json github_secrets  # 機密値（-json が必要）

# ── GitHub Variables（Settings > Variables）────────────
# 非機密の設定値。ワークフローから ${{ vars.XXX }} で参照する

output "github_variables" {
  description = "GitHub Actions の Variables に設定する値（非機密）"
  value = {
    GCP_PROJECT_ID      = var.gcp_project_id
    GCP_REGION          = var.gcp_region
    AR_REPO             = google_artifact_registry_repository.backend.repository_id
    CLOUD_RUN_SERVICE   = google_cloud_run_v2_service.backend.name
    WIF_PROVIDER        = google_iam_workload_identity_pool_provider.github.name
    WIF_SERVICE_ACCOUNT = google_service_account.deployer.email
  }
}

# ── GitHub Secrets（Settings > Secrets）───────────────
# 機密値。ワークフローから ${{ secrets.XXX }} で参照する

output "github_secrets" {
  description = "GitHub Actions の Secrets に設定する値（機密）"
  sensitive   = true
  value = {
    DATABASE_URL = local.database_url
  }
}

# ── 参考情報 ───────────────────────────────────────────

output "cloud_run_url" {
  description = "バックエンドの公開 URL"
  value       = "https://${google_cloud_run_v2_service.backend.uri}"
}

output "vercel_project_id" {
  description = "Vercel プロジェクト ID"
  value       = vercel_project.frontend.id
}

output "artifact_registry_hostname" {
  description = "Docker イメージの push 先ホスト名"
  value       = "${var.gcp_region}-docker.pkg.dev"
}
