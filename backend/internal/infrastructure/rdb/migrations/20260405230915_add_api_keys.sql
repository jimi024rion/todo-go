-- Modify "todos" table
ALTER TABLE "public"."todos" DROP CONSTRAINT "chk_todos_status", ADD CONSTRAINT "chk_todos_status" CHECK ((status)::text = ANY ((ARRAY['pending'::character varying, 'in_progress'::character varying, 'completed'::character varying])::text[]));
-- Create "api_keys" table
CREATE TABLE "public"."api_keys" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "user_id" uuid NOT NULL, "key" character varying(64) NOT NULL, "name" character varying(255) NOT NULL, "created_at" timestamptz NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "api_keys_key_key" UNIQUE ("key"), CONSTRAINT "fk_api_keys_user_id" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
-- Set comment to table: "api_keys"
COMMENT ON TABLE "public"."api_keys" IS 'APIキー';
