# ── GCP ─────────────────────────────────────────────────

variable "gcp_project_id" {
  description = "GCP プロジェクト ID"
  type        = string
}

variable "gcp_region" {
  description = "Cloud Run などを配置するリージョン"
  type        = string
  default     = "asia-northeast1"
}

# ── GitHub ───────────────────────────────────────────────

variable "github_repo" {
  description = "GitHub リポジトリ（owner/repo 形式）例: jimi024rion/todo-go"
  type        = string
}

# ── Neon ─────────────────────────────────────────────────

variable "neon_api_key" {
  description = "Neon API キー（https://console.neon.tech/app/settings/api-keys で発行）"
  type        = string
  sensitive   = true
}

variable "neon_org_id" {
  description = "Neon Organization ID（https://console.neon.tech/app/settings/organization で確認）"
  type        = string
}

