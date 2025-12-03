-- Create "todos" table
CREATE TABLE "public"."todos" ("id" serial NOT NULL, "title" character varying(255) NOT NULL, "description" text NULL, "completed" boolean NOT NULL DEFAULT false, "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY ("id"));
-- Drop "users" table
DROP TABLE "public"."users";
