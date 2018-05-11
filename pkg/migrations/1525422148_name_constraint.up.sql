ALTER TABLE solutions
  ALTER COLUMN name DROP NOT NULL,
  DROP CONSTRAINT "solutions_name_key",
  ADD CONSTRAINT "name_user_id" UNIQUE( "name", "user_id" );