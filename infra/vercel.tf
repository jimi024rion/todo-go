# Vercel は現在スキップ。フロントエンドをデプロイする際にコメントを解除する。
# versions.tf / main.tf / variables.tf の vercel ブロックも同時に復活させること。
#
# resource "vercel_project" "frontend" {
#   name      = "todo-go-frontend"
#   framework = "nextjs"
#   root_directory = "frontend"
#   git_repository = {
#     type = "github"
#     repo = var.github_repo
#   }
# }
#
# resource "vercel_project_environment_variable" "backend_url_production" {
#   project_id = vercel_project.frontend.id
#   key        = "BACKEND_URL"
#   value      = "https://${google_cloud_run_v2_service.backend.uri}"
#   target     = ["production"]
# }
#
# resource "vercel_project_environment_variable" "backend_url_preview" {
#   project_id = vercel_project.frontend.id
#   key        = "BACKEND_URL"
#   value  = "http://localhost:8080"
#   target = ["preview"]
# }
