-- Modify "api_keys" table
ALTER TABLE "public"."api_keys" DROP COLUMN "key_prefix", ALTER COLUMN "key_hash" TYPE character varying(64), ADD CONSTRAINT "api_keys_key_hash_key" UNIQUE ("key_hash");
-- Create index "idx_api_keys_key_hash" to table: "api_keys"
CREATE INDEX "idx_api_keys_key_hash" ON "public"."api_keys" ("key_hash");
-- Modify "todos" table
ALTER TABLE "public"."todos" DROP CONSTRAINT "chk_todos_status", ADD CONSTRAINT "chk_todos_status" CHECK ((status)::text = ANY ((ARRAY['pending'::character varying, 'in_progress'::character varying, 'completed'::character varying])::text[]));
