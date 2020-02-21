DROP TRIGGER update_candidates_updated_at on candidates;

ALTER TABLE candidates DROP COLUMN created_at;
ALTER TABLE candidates DROP COLUMN updated_at;
