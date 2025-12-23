CREATE TABLE "users" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "name" varchar(255) NOT NULL,
  "email" varchar(255) NOT NULL UNIQUE,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL
);

CREATE TABLE "todos" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL CONSTRAINT "fk_todos_user_id" REFERENCES "users" ("id") ON DELETE CASCADE,
  "title" varchar(255) NOT NULL,
  "description" text,
  "status" varchar(50) NOT NULL DEFAULT 'pending' CONSTRAINT "chk_todos_status" CHECK (status IN ('pending', 'in_progress', 'completed')),
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL
);


CREATE TABLE "tags" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL CONSTRAINT "fk_tags_user_id" REFERENCES "users" ("id") ON DELETE CASCADE,
  "name" varchar(255) NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL
);

CREATE TABLE "todo_tags" (
  "todo_id" uuid NOT NULL CONSTRAINT "fk_todo_tags_todo_id" REFERENCES "todos" ("id") ON DELETE CASCADE,
  "tag_id" uuid NOT NULL CONSTRAINT "fk_todo_tags_tag_id" REFERENCES "tags" ("id") ON DELETE CASCADE,
  PRIMARY KEY ("todo_id", "tag_id")
);

COMMENT ON TABLE "users" IS 'ユーザー情報';
COMMENT ON TABLE "todos" IS 'ToDoタスク';
COMMENT ON COLUMN "todos"."status" IS 'タスクの状態 (pending, in_progress, completed)';
COMMENT ON TABLE "tags" IS 'タグ';
COMMENT ON TABLE "todo_tags" IS 'ToDoとタグの中間テーブル';
