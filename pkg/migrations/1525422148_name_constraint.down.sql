ALTER TABLE solutions
  ALTER COLUMN name SET NOT NULL;
ALTER TABLE solutions
  DROP CONSTRAINT "name_user_id";