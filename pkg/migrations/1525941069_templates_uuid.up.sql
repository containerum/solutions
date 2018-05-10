ALTER TABLE templates
  DROP CONSTRAINT available_solutions_pkey,
  ADD COLUMN id UUID DEFAULT uuid_generate_v4() PRIMARY KEY NOT NULL;
CREATE UNIQUE INDEX only_active_template
  ON templates (name) WHERE (active);
ALTER TABLE solutions
  ADD COLUMN template_id UUID NOT NULL,
  DROP COLUMN template,
  ADD CONSTRAINT solutions_template_fkey FOREIGN KEY (template_id) REFERENCES templates (id),
  ALTER COLUMN is_deleted SET DEFAULT 'false',
  ALTER COLUMN is_deleted SET NOT NULL,
  DROP CONSTRAINT "name_user_id";
CREATE UNIQUE INDEX only_not_deleted_solution
  ON solutions (name, user_id) WHERE (NOT is_deleted);
