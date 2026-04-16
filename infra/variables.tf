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

# ── Vercel ───────────────────────────────────────────────

variable "vercel_token" {
  description = "Vercel API トークン（https://vercel.com/account/tokens で発行）"
  type        = string
  sensitive   = true
}

variable "vercel_team_id" {
  description = "Vercel Team ID（個人アカウントの場合は空文字列）"
  type        = string
  default     = ""
}
