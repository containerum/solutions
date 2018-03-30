CREATE TABLE IF NOT EXISTS available_solutions
(
  name TEXT PRIMARY KEY NOT NULL,
  cpu TEXT NOT NULL,
  ram TEXT NOT NULL,
  images jsonb NOT NULL,
  url text NOT NULL
);