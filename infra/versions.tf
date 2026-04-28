terraform {
  required_version = ">= 1.9.0"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 6.0"
    }
    neon = {
      source  = "kislerdm/neon"
      version = "~> 0.9"
    }
  }

  # ── リモートバックエンド（推奨） ──────────────────────────
  # ローカルの terraform.tfstate は git に含めないこと。
  # チーム開発・CI/CD での apply には GCS バックエンドを推奨。
  #
  # 事前にバケットを作成してから有効化する:
  #   gcloud storage buckets create gs://YOUR_BUCKET --location=asia-northeast1
  #
  # backend "gcs" {
  #   bucket = "YOUR_BUCKET"
  #   prefix = "todo-go/state"
  # }
  # ────────────────────────────────────────────────────────
}
