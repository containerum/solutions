CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS solutions
(
  id UUID DEFAULT uuid_generate_v4() PRIMARY KEY NOT NULL,
  template TEXT NOT NULL,
  name TEXT UNIQUE NOT NULL,
  namespace TEXT NOT NULL,
  user_id UUID NOT NULL,
  created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT now() NOT NULL,
  is_deleted BOOLEAN,
  deleted_at TIMESTAMP WITHOUT TIME ZONE
);
CREATE TABLE IF NOT EXISTS parameters
(
  id UUID DEFAULT uuid_generate_v4() PRIMARY KEY NOT NULL,
  solution_id UUID NOT NULL,
  branch TEXT NOT NULL,
  env jsonb,
  data jsonb,
  CONSTRAINT solutions_parametes_fkey FOREIGN KEY (solution_id) REFERENCES solutions (id) ON DELETE CASCADE
);
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