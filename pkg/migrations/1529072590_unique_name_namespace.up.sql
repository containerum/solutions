DROP INDEX only_not_deleted_solution;
CREATE UNIQUE INDEX only_not_deleted_solution_namespace
  ON solutions (name, namespace) WHERE (NOT is_deleted);
