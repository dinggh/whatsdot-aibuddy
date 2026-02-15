ALTER TABLE homework_records
  ALTER COLUMN user_id DROP NOT NULL;

ALTER TABLE homework_records
  ADD COLUMN IF NOT EXISTS device_id TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS mode TEXT NOT NULL DEFAULT 'guided',
  ADD COLUMN IF NOT EXISTS summary TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS question_text TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS result_json JSONB NOT NULL DEFAULT '{}'::jsonb,
  ADD COLUMN IF NOT EXISTS source_image_url TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

CREATE INDEX IF NOT EXISTS idx_homework_records_device_id_solved_at
  ON homework_records(device_id, solved_at DESC);
