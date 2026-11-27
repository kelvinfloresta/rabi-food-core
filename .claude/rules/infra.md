---
paths:
  - "infra/**"
  - ".github/**"
---

# Infrastructure

## Checks
- `task validate_alloy` for `infra/monitoring/config.alloy`; `terraform fmt -check` and `terraform validate` in `infra/terraform`. Known validate failure: `stripe_apple_pay_domain` in `stripe.tf` is not supported by the pinned Stripe provider.
- Terraform is applied by hand (`task tf:apply`), never by CI; prepare changes and let the user run `plan`/`apply`.

## Logs and metrics
- One Alloy config, `infra/monitoring/config.alloy`, runs in production (ECS sidecar) and in the local test (`task observe`). Differences go in environment variables, never in a second config: what is tested locally is what runs in production.
- Containers whose logs matter log to Docker (`json-file`, `local.docker_log_configuration`), depend on `alloy` with condition `START` (it starts first and stops last, to read a crash's final lines; never `HEALTHY`, which would let a broken Alloy keep the API down) and are listed in Alloy's `LOG_SERVICES`. Do not give them `awslogs`: CloudWatch is only for what Alloy cannot read (Alloy itself, the one-shot `gotrue-migration`).
- People's data never leaves for Grafana Cloud (decision 010): a service whose lines carry e-mails, user ids or IPs gets a masking `stage.match` in `loki.process`, as GoTrue has; test new patterns against sample lines with `loki.echo` first.
- Accepted overlap: while a deploy runs the old and new task side by side on one instance, both Alloys ship both tasks' lines and share `/var/lib/alloy`. Loki drops exact duplicates; the bytes are still billed. Moving Alloy to one daemon per instance would end it.
- Besides `service` and `env`, only `level` becomes a Loki label. Request fields stay in the JSON line and are parsed at query time — a new label multiplies the index.
- Keep the Alloy image version the same in `infra/terraform/ecs.tf`, `infra/monitoring/docker-compose.yaml` and `Taskfile.yml`.

## Platform dashboards
- Grafana Cloud reaches the private database through the `pdc-agent` container; never make RDS public or open its security group to reach it.
- Dashboards live in `infra/terraform/dashboards/*.json` and query the `analytics` views only, with plain SQL (no Grafana macros): the analytics integration test runs each query as `analytics_reader`.
- The `grafana_reader` login is created by `deploy.yml` from the `rabi-food/grafana-reader-password` secret; no password goes in a migration. Its `CONNECTION LIMIT` and the data source's `maxOpenConns` (`grafana.tf`) move together: the pool must stay below the limit, or panels fail at random. The CI user (not managed by Terraform) needs `secretsmanager:GetSecretValue` on that secret, and the secret must exist (`terraform apply`) before the release.

## ECS task definition
- Releases register a new revision from the live one and edit `containerDefinitions[0]` (`.github/workflows/deploy.yml`): `api` must stay the first container.
- `ignore_changes = [container_definitions]`: a changed container definition only takes effect with `terraform apply -replace=aws_ecs_task_definition.api` (a change to `config.alloy` triggers it through `terraform_data.alloy_config`), and the running service only switches on the next release. Run a release right after the apply.
- Removing a log group or other resource the running revision still uses: use a `removed` block with `destroy = false`, never a plain delete — the old revision keeps running until the next release.
- A new secret in a container needs its ARN in the execution role's `secretsmanager:GetSecretValue` policy (`iam.tf`), or the task fails to start.
