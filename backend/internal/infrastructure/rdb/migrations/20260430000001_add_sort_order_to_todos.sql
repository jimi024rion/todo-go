-- Add sort_order column for user-defined ordering
ALTER TABLE "todos" ADD COLUMN "sort_order" FLOAT NOT NULL DEFAULT 0;

-- Backfill existing rows: use negative unix timestamp so newer items sort first
UPDATE "todos" SET "sort_order" = EXTRACT(EPOCH FROM "created_at") * -1;
