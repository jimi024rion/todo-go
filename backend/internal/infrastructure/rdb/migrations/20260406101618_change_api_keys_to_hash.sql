-- Modify "api_keys" table
ALTER TABLE "public"."api_keys" DROP COLUMN "key", ADD COLUMN "key_prefix" character varying(12) NOT NULL, ADD COLUMN "key_hash" character varying(72) NOT NULL;
-- Create index "idx_api_keys_key_prefix" to table: "api_keys"
CREATE INDEX "idx_api_keys_key_prefix" ON "public"."api_keys" ("key_prefix");
-- Modify "todos" table
ALTER TABLE "public"."todos" DROP CONSTRAINT "chk_todos_status", ADD CONSTRAINT "chk_todos_status" CHECK ((status)::text = ANY ((ARRAY['pending'::character varying, 'in_progress'::character varying, 'completed'::character varying])::text[]));
