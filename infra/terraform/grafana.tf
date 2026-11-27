data "grafana_cloud_stack" "main" {
  slug = var.grafana_stack_slug
}

# Grafana Cloud reaches the private database through the PDC agent in the API task (ecs.tf): RDS stays private.
resource "grafana_cloud_private_data_source_connect_network" "main" {
  region           = data.grafana_cloud_stack.main.region_slug
  name             = "rabi-food-vpc"
  display_name     = "Rabi Food VPC"
  stack_identifier = data.grafana_cloud_stack.main.id
}

resource "grafana_cloud_private_data_source_connect_network_token" "agent" {
  pdc_network_id = grafana_cloud_private_data_source_connect_network.main.pdc_network_id
  region         = grafana_cloud_private_data_source_connect_network.main.region
  name           = "rabi-food-pdc-agent"
  display_name   = "Rabi Food PDC agent"
}

# Platform metrics: the analytics views only (docs/decisions/010-platform-business-metrics.md).
resource "grafana_data_source" "analytics" {
  type     = "grafana-postgresql-datasource"
  name     = "Rabi Food analytics"
  uid      = "rabi-food-analytics"
  url      = "${aws_db_instance.postgres.address}:5432"
  username = "grafana_reader"

  private_data_source_connect_network_id = grafana_cloud_private_data_source_connect_network.main.pdc_network_id

  # The pool stays under grafana_reader's CONNECTION LIMIT (deploy.yml), leaving one connection for a manual psql.
  json_data_encoded = jsonencode({
    database        = var.db_name
    sslmode         = "require"
    maxOpenConns    = 2
    maxIdleConns    = 2
    connMaxLifetime = 14400
  })

  secure_json_data_encoded = jsonencode({
    password = var.grafana_reader_password
  })
}

resource "grafana_dashboard" "platform" {
  config_json = file("${path.module}/dashboards/platform.json")
}
