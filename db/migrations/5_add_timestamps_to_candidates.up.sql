ALTER TABLE candidates
  ADD COLUMN created_at TIMESTAMP DEFAULT now(),
  ADD COLUMN updated_at TIMESTAMP DEFAULT now(); 

CREATE TRIGGER update_candidates_updated_at 
  BEFORE UPDATE ON candidates
  FOR EACH ROW EXECUTE PROCEDURE modify_updated_at_column();
