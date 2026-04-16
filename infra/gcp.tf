# ═══════════════════════════════════════════════════════
# 1. API 有効化
# ═══════════════════════════════════════════════════════
# GCP のリソースを使う前に該当サービスの API を有効化する必要がある
# disable_on_destroy = false にすることで terraform destroy 時に API を無効化しない
# （他のリソースが依存している可能性があるため）

resource "google_project_service" "run" {
  service            = "run.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "artifactregistry" {
  service            = "artifactregistry.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "secretmanager" {
  service            = "secretmanager.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "iam" {
  service            = "iam.googleapis.com"
  disable_on_destroy = false
}

# ═══════════════════════════════════════════════════════
# 2. Artifact Registry（Docker イメージ保管庫）
# ═══════════════════════════════════════════════════════

resource "google_artifact_registry_repository" "backend" {
  location      = var.gcp_region
  repository_id = "todo-go"
  format        = "DOCKER"
  description   = "todo-go バックエンド Docker イメージ"

  depends_on = [google_project_service.artifactregistry]
}

# ═══════════════════════════════════════════════════════
# 3. サービスアカウント（GitHub Actions がデプロイに使用）
# ═══════════════════════════════════════════════════════

resource "google_service_account" "deployer" {
  account_id   = "deployer"
  display_name = "GitHub Actions Deployer"
  description  = "GitHub Actions の CD ワークフローがデプロイに使用するサービスアカウント"

  depends_on = [google_project_service.iam]
}

# ── IAM ロール付与 ────────────────────────────────────────
# for_each でまとめて付与する（Go の range に相当）

locals {
  deployer_roles = toset([
    "roles/run.admin",               # Cloud Run サービスの作成・更新
    "roles/artifactregistry.writer", # Docker イメージの push
    "roles/secretmanager.admin",     # Secret Manager の管理
    "roles/iam.serviceAccountUser",  # サービスアカウントの権限借用
  ])
}

resource "google_project_iam_member" "deployer" {
  for_each = local.deployer_roles

  project = var.gcp_project_id
  role    = each.value
  member  = "serviceAccount:${google_service_account.deployer.email}"
}

# ═══════════════════════════════════════════════════════
# 4. Workload Identity Federation（WIF）
# ═══════════════════════════════════════════════════════
# GitHub Actions がサービスアカウントキー（JSON）なしで GCP に認証できる仕組み
# GitHub の OIDC トークンを GCP のアイデンティティとして信頼する

resource "google_iam_workload_identity_pool" "github" {
  workload_identity_pool_id = "github-pool"
  display_name              = "GitHub Actions Pool"
  description               = "GitHub Actions からのデプロイ用 WIF プール"

  depends_on = [google_project_service.iam]
}

resource "google_iam_workload_identity_pool_provider" "github" {
  workload_identity_pool_id          = google_iam_workload_identity_pool.github.workload_identity_pool_id
  workload_identity_pool_provider_id = "github-provider"
  display_name                       = "GitHub Actions Provider"

  # GitHub の OIDC トークンから GCP のアイデンティティへの変換ルール
  attribute_mapping = {
    "google.subject"       = "assertion.sub"        # GitHub Actions の実行識別子
    "attribute.repository" = "assertion.repository" # リポジトリ名（owner/repo）
  }

  # このリポジトリからのリクエストのみ許可
  attribute_condition = "assertion.repository == '${var.github_repo}'"

  oidc {
    issuer_uri = "https://token.actions.githubusercontent.com"
  }
}

# deployer SA が GitHub Actions（特定リポジトリ）から呼び出せるようにする
resource "google_service_account_iam_member" "wif_deployer" {
  service_account_id = google_service_account.deployer.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.github.name}/attribute.repository/${var.github_repo}"
}

# ═══════════════════════════════════════════════════════
# 5. Secret Manager（DB 認証情報の安全な保管）
# ═══════════════════════════════════════════════════════
# DB の認証情報を Secret Manager に保管し、Cloud Run がマウントして使う
# → 環境変数に平文で書かない、GitHub Secrets にも入れない

resource "google_secret_manager_secret" "db" {
  for_each = local.db_secrets

  secret_id = each.key
  replication {
    auto {} # GCP が自動でレプリケーション先を選択
  }

  depends_on = [google_project_service.secretmanager]
}

resource "google_secret_manager_secret_version" "db" {
  for_each = local.db_secrets

  secret      = google_secret_manager_secret.db[each.key].id
  secret_data = each.value # Neon の出力値を参照
}

# Cloud Run のデフォルト SA が Secret Manager を読み取れるようにする
resource "google_secret_manager_secret_iam_member" "cloudrun_accessor" {
  for_each = local.db_secrets

  secret_id = google_secret_manager_secret.db[each.key].secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${local.cloudrun_sa}"

  depends_on = [google_secret_manager_secret.db]
}

# ═══════════════════════════════════════════════════════
# 6. Cloud Run サービス
# ═══════════════════════════════════════════════════════

resource "google_cloud_run_v2_service" "backend" {
  name     = "todo-go-backend"
  location = var.gcp_region

  template {
    containers {
      # 初回は Cloud Run の公式プレースホルダーイメージを使用
      # 以降は CI/CD（cd.yml）が実際のイメージに更新する
      image = "us-docker.pkg.dev/cloudrun/container/hello"

      ports {
        container_port = 8080
      }

      # 非機密の環境変数
      env {
        name  = "APP_ENV"
        value = "production"
      }
      env {
        name  = "PORT"
        value = "8080"
      }
      env {
        name  = "LOG_LEVEL"
        value = "info"
      }
      env {
        name  = "DB_PORT"
        value = "5432"
      }

      # Secret Manager からマウントする環境変数
      dynamic "env" {
        for_each = local.db_secrets
        content {
          name = env.key
          value_source {
            secret_key_ref {
              secret  = google_secret_manager_secret.db[env.key].secret_id
              version = "latest"
            }
          }
        }
      }
    }
  }

  lifecycle {
    # CI/CD がイメージを更新するたびに Terraform との差分が出るため無視する
    # Goで言う「このフィールドは外部から変更される」という宣言
    ignore_changes = [
      template[0].containers[0].image,
    ]
  }

  depends_on = [
    google_project_service.run,
    google_secret_manager_secret_version.db,
  ]
}

# Cloud Run サービスをインターネットに公開（認証不要）
resource "google_cloud_run_v2_service_iam_member" "public" {
  location = google_cloud_run_v2_service.backend.location
  name     = google_cloud_run_v2_service.backend.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}
