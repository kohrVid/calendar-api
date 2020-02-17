CREATE UNIQUE INDEX candidates_unique_idx
  ON candidates
  USING btree (first_name, last_name, email);
