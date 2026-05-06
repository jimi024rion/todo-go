-- Modify "todos" table
ALTER TABLE "public"."todos" DROP CONSTRAINT "chk_todos_status", ADD CONSTRAINT "chk_todos_status" CHECK ((status)::text = ANY ((ARRAY['pending'::character varying, 'in_progress'::character varying, 'completed'::character varying])::text[])), ALTER COLUMN "sort_order" SET DEFAULT 0.0;
-- Modify "users" table
ALTER TABLE "public"."users" ADD COLUMN "firebase_uid" character varying(128) NULL, ADD CONSTRAINT "users_firebase_uid_key" UNIQUE ("firebase_uid");
