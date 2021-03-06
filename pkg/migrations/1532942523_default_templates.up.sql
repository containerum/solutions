INSERT INTO templates VALUES ('mariadb-solution', '100m', '128Mi', '["mariadb"]', 'https://github.com/containerum/mariadb-solution', true) ON CONFLICT DO NOTHING;
INSERT INTO templates VALUES ('postgresql-solution', '100m', '128Mi', '["postgres"]', 'https://github.com/containerum/postgresql-solution', true) ON CONFLICT DO NOTHING;
INSERT INTO templates VALUES ('grafana-xxl-solution', '100m', '128Mi', '["grafana-xxl"]', 'https://github.com/containerum/grafana-xxl-solution', true) ON CONFLICT DO NOTHING;
INSERT INTO templates VALUES ('rabbitmq-manager-solution', '100m', '128Mi', '["rabbitmq"]', 'https://github.com/containerum/rabbitmq-manager-solution', true) ON CONFLICT DO NOTHING;
INSERT INTO templates VALUES ('webpack-3.8-ssh-solution', '100m', '128Mi', '["webpack-3.8-ssh-solution"]', 'https://github.com/containerum/webpack-3.8-ssh-solution', true) ON CONFLICT DO NOTHING;
INSERT INTO templates VALUES ('redmine-mariadb-solution', '200m', '256Mi', '["redmine", "mariadb"]', 'https://github.com/containerum/redmine-mariadb-solution', true) ON CONFLICT DO NOTHING;
INSERT INTO templates VALUES ('redmine-postgresql-solution', '200m', '256Mi', '["redmine", "postgres"]', 'https://github.com/containerum/redmine-postgresql-solution', true) ON CONFLICT DO NOTHING;
INSERT INTO templates VALUES ('magento-solution', '200m', '256Mi', '["mysql", "magento2"]', 'https://github.com/containerum/magento-solution', true) ON CONFLICT DO NOTHING;
