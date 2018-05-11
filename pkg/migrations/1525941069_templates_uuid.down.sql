ALTER TABLE templates
  DROP COLUMN id,
  ADD PRIMARY KEY (name);
DROP INDEX only_active_template;
ALTER TABLE solutions
  DROP COLUMN template_id,
  ADD COLUMN template TEXT NOT NULL,
  DROP CONSTRAINT solutions_template_fkey,
  ALTER COLUMN is_deleted DROP DEFAULT,
  ALTER COLUMN is_deleted DROP NOT NULL,
  ADD CONSTRAINT "name_user_id" UNIQUE( "name", "user_id" );
DROP INDEX only_not_deleted_solution;