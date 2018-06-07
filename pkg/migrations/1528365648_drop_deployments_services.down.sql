CREATE TABLE IF NOT EXISTS deployments
(
  id UUID DEFAULT uuid_generate_v4() PRIMARY KEY NOT NULL,
  deploy_name TEXT NOT NULL,
  solution_id UUID NOT NULL,
  CONSTRAINT deployments_solutions_fkey FOREIGN KEY (solution_id) REFERENCES solutions (id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS services
(
  id UUID DEFAULT uuid_generate_v4() PRIMARY KEY NOT NULL,
  service_name TEXT NOT NULL,
  solution_id UUID NOT NULL,
  CONSTRAINT services_solutions_fkey FOREIGN KEY (solution_id) REFERENCES solutions (id) ON DELETE CASCADE
);