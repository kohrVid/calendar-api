CREATE TABLE time_slots (
  id SERIAL,
  start_time INT NOT NULL,
  duration INT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW() NOT NULL,
  updated_at TIMESTAMP DEFAULT NOW() NOT NULL
);

CREATE TRIGGER update_time_slots_updated_at
  BEFORE UPDATE ON time_slots
  FOR EACH ROW EXECUTE PROCEDURE modify_updated_at_column();
