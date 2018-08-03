DROP INDEX only_not_deleted_solution_namespace;
CREATE UNIQUE INDEX only_not_deleted_solution
  ON solutions (name, user_id) WHERE (NOT is_deleted);
