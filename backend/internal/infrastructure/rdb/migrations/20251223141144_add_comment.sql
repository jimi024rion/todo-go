-- Set comment to table: "users"
COMMENT ON TABLE "public"."users" IS 'ユーザー情報';
-- Modify "tags" table
ALTER TABLE "public"."tags" DROP CONSTRAINT "tags_user_id_fkey", ADD CONSTRAINT "fk_tags_user_id" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Set comment to table: "tags"
COMMENT ON TABLE "public"."tags" IS 'タグ';
-- Modify "todos" table
ALTER TABLE "public"."todos" DROP CONSTRAINT "todos_status_check", DROP CONSTRAINT "todos_user_id_fkey", ADD CONSTRAINT "chk_todos_status" CHECK ((status)::text = ANY ((ARRAY['pending'::character varying, 'in_progress'::character varying, 'completed'::character varying])::text[])), ADD CONSTRAINT "fk_todos_user_id" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Set comment to table: "todos"
COMMENT ON TABLE "public"."todos" IS 'ToDoタスク';
-- Set comment to column: "status" on table: "todos"
COMMENT ON COLUMN "public"."todos"."status" IS 'タスクの状態 (pending, in_progress, completed)';
-- Modify "todo_tags" table
ALTER TABLE "public"."todo_tags" DROP CONSTRAINT "todo_tags_tag_id_fkey", DROP CONSTRAINT "todo_tags_todo_id_fkey", ADD CONSTRAINT "fk_todo_tags_tag_id" FOREIGN KEY ("tag_id") REFERENCES "public"."tags" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT "fk_todo_tags_todo_id" FOREIGN KEY ("todo_id") REFERENCES "public"."todos" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Set comment to table: "todo_tags"
COMMENT ON TABLE "public"."todo_tags" IS 'ToDoとタグの中間テーブル';
