-- Modify "todos" table
ALTER TABLE "todos" ADD CONSTRAINT "todos_status_check" CHECK ((status)::text = ANY ((ARRAY['pending'::character varying, 'in_progress'::character varying, 'completed'::character varying])::text[])), DROP COLUMN "done", ADD COLUMN "status" character varying(50) NOT NULL DEFAULT 'pending';
