ALTER TABLE solutions
  ALTER COLUMN name DROP NOT NULL;
ALTER TABLE solutions
  DROP CONSTRAINT "solutions_name_key";
ALTER TABLE solutions
  ADD CONSTRAINT "name_user_id" UNIQUE( "name", "user_id" );