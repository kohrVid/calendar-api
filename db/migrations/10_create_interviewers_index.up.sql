CREATE UNIQUE INDEX interviewers_unique_idx
  ON interviewers
  USING btree (first_name, last_name, email);
