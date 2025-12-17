// Define an environment named "local"
env "local" {
  // Declare where the schema definition resides.
  // Also supported: ["file://multi.hcl", "file://schema.hcl"].
  src = "file://schema.hcl"

  // Define the URL of the database which is managed
  // in this environment.
  url = "postgres://user:P@ssw0rd@localhost:5432/todo-go?search_path=public&sslmode=disable"

  // Define the URL of the Dev Database for this environment
  // See: https://atlasgo.io/concepts/dev-database
  dev = "docker://postgres/15/dev?search_path=public"
}

env "dev" {
  // ... a different env
}
