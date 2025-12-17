# migration

マイグレーション時のメモ

初回

1. [atlas.hcl](./atlas.hcl)作成
2. [schema.sql](./schema.sql)作成
3. マイグレーションファイル生成

    ```shell
    ❯ atlas migrate diff initial \
        --to file://schema.sql \
        --dev-url "docker://postgres/15/dev?search_path=public" \
        --format '{{ sql . "  " }}'
    ```

4. `migrations`ディレクトリ直下に以下ファイルが生成される。

   - マイグレーション用のSQL: [20251217135055_initial.sql](./migrations/20251217135055_initial.sql)
   - チェックサム: [./migrations/atlas.sum](./migrations/atlas.sum)

5. マイグレーション実行

    ```shell
    ❯ atlas migrate apply --env local
    Migrating to version 20251217135055 (1 migrations in total):

    -- migrating version 20251217135055
        -> CREATE TABLE "users" (
            "id" uuid NOT NULL DEFAULT gen_random_uuid(),
            "name" character varying(255) NOT NULL,
            "email" character varying(255) NOT NULL,
            "created_at" timestamptz NOT NULL,
            "updated_at" timestamptz NOT NULL,
            PRIMARY KEY ("id"),
            CONSTRAINT "users_email_key" UNIQUE ("email")
        );
        -> CREATE TABLE "tags" (
            "id" uuid NOT NULL DEFAULT gen_random_uuid(),
            "user_id" uuid NOT NULL,
            "name" character varying(255) NOT NULL,
            "created_at" timestamptz NOT NULL,
            "updated_at" timestamptz NOT NULL,
            PRIMARY KEY ("id"),
            CONSTRAINT "tags_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
        );
        -> CREATE TABLE "todos" (
            "id" uuid NOT NULL DEFAULT gen_random_uuid(),
            "user_id" uuid NOT NULL,
            "title" character varying(255) NOT NULL,
            "description" text NULL,
            "done" boolean NOT NULL DEFAULT false,
            "created_at" timestamptz NOT NULL,
            "updated_at" timestamptz NOT NULL,
            PRIMARY KEY ("id"),
            CONSTRAINT "todos_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
        );
        -> CREATE TABLE "todo_tags" (
            "todo_id" uuid NOT NULL,
            "tag_id" uuid NOT NULL,
            PRIMARY KEY ("todo_id", "tag_id"),
            CONSTRAINT "todo_tags_tag_id_fkey" FOREIGN KEY ("tag_id") REFERENCES "tags" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
            CONSTRAINT "todo_tags_todo_id_fkey" FOREIGN KEY ("todo_id") REFERENCES "todos" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
        );
    -- ok (16.916625ms)

    -------------------------
    -- 60.656042ms
    -- 1 migration
    -- 4 sql statements
    ```
