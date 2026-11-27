-- +goose Up
-- Platform metrics read by Grafana (docs/decisions/010-platform-business-metrics.md). The views are the contract the
-- dashboards depend on: aggregates per tenant only, never a person's data. Months are counted in Brazil's time.
CREATE SCHEMA IF NOT EXISTS analytics;

-- A sale is an order or a subscription; paid_cents is set only once it was paid. A subscription is one row per
-- product, so sale_key groups the rows of one checkout: a sale is what the customer bought, not a product. It
-- carries user_id to count buyers, so it is not granted to the reader: only the per-tenant views below are.
CREATE VIEW
    analytics.sales AS
SELECT
    tenant_id,
    user_id,
    id::TEXT AS sale_key,
    FALSE AS renewal,
    'order' AS kind,
    created_at,
    paid_at,
    CASE
        WHEN payment_status = 'paid' THEN total_price
    END AS paid_cents
FROM
    commerce.orders
WHERE
    deleted_at IS NULL
UNION ALL
SELECT
    tenant_id,
    user_id,
    COALESCE(checkout_key, id)::TEXT AS sale_key,
    root_subscription_id IS NOT NULL AS renewal,
    'subscription' AS kind,
    created_at,
    paid_at,
    CASE
        WHEN payment_status = 'paid' THEN payment_amount
    END AS paid_cents
FROM
    SUBSCRIPTION.subscriptions;

-- The owner's membership is written in the same transaction as the tenant, so its date is the tenant's.
CREATE VIEW
    analytics.tenants AS
SELECT
    t.id AS tenant_id,
    t.name,
    t.slug,
    MIN(m.created_at) AS created_at,
    COUNT(m.user_id) AS members
FROM
    iam.tenants t
    LEFT JOIN iam.tenant_members m ON m.tenant_id = t.id
GROUP BY
    t.id;

CREATE VIEW
    analytics.tenant_monthly AS
WITH
    events AS (
        SELECT
            tenant_id,
            user_id,
            sale_key,
            renewal,
            kind,
            date_trunc('month', created_at AT TIME ZONE 'America/Sao_Paulo')::date AS MONTH,
            FALSE AS paid,
            0::BIGINT AS paid_cents
        FROM
            analytics.sales
        WHERE
            created_at IS NOT NULL
        UNION ALL
        SELECT
            tenant_id,
            user_id,
            sale_key,
            renewal,
            kind,
            date_trunc('month', paid_at AT TIME ZONE 'America/Sao_Paulo')::date,
            TRUE,
            paid_cents
        FROM
            analytics.sales
        WHERE
            paid_at IS NOT NULL
            AND paid_cents IS NOT NULL
    ),
    first_purchase AS (
        SELECT
            tenant_id,
            user_id,
            MIN(MONTH) AS MONTH
        FROM
            events
        WHERE
            paid
        GROUP BY
            tenant_id,
            user_id
    )
SELECT
    e.tenant_id,
    e.month,
    COUNT(DISTINCT e.user_id) FILTER (
        WHERE
            e.paid
    ) AS buyers,
    (
        SELECT
            COUNT(*)
        FROM
            first_purchase f
        WHERE
            f.tenant_id = e.tenant_id
            AND f.month = e.month
    ) AS new_buyers,
    COUNT(DISTINCT e.sale_key) FILTER (
        WHERE
            e.kind = 'order'
            AND NOT e.paid
    ) AS orders_created,
    COUNT(DISTINCT e.sale_key) FILTER (
        WHERE
            e.kind = 'order'
            AND e.paid
    ) AS orders_paid,
    COALESCE(
        SUM(e.paid_cents) FILTER (
            WHERE
                e.kind = 'order'
                AND e.paid
        ),
        0
    )::BIGINT AS orders_paid_cents,
    COUNT(DISTINCT e.sale_key) FILTER (
        WHERE
            e.kind = 'subscription'
            AND NOT e.paid
            AND NOT e.renewal
    ) AS subscriptions_created,
    COUNT(DISTINCT e.sale_key) FILTER (
        WHERE
            e.kind = 'subscription'
            AND NOT e.paid
            AND e.renewal
    ) AS subscription_renewals,
    COUNT(DISTINCT e.sale_key) FILTER (
        WHERE
            e.kind = 'subscription'
            AND e.paid
    ) AS subscriptions_paid,
    COALESCE(
        SUM(e.paid_cents) FILTER (
            WHERE
                e.kind = 'subscription'
                AND e.paid
        ),
        0
    )::BIGINT AS subscriptions_paid_cents
FROM
    events e
GROUP BY
    e.tenant_id,
    e.month;

CREATE VIEW
    analytics.tenant_current AS
SELECT
    t.id AS tenant_id,
    (
        SELECT
            COUNT(DISTINCT COALESCE(s.checkout_key, s.id))
        FROM
            SUBSCRIPTION.subscriptions s
        WHERE
            s.tenant_id = t.id
            AND s.status = 'active'
            AND s.payment_status = 'paid'
    ) AS active_subscriptions,
    (
        SELECT
            COUNT(DISTINCT s.user_id)
        FROM
            analytics.sales s
        WHERE
            s.tenant_id = t.id
            AND s.paid_cents IS NOT NULL
    ) AS buyers,
    (
        SELECT
            MAX(s.created_at)
        FROM
            analytics.sales s
        WHERE
            s.tenant_id = t.id
    ) AS last_activity_at
FROM
    iam.tenants t;

-- Grafana's login role is granted analytics_reader outside the migrations, so no password lives here.
-- Views run with their owner's rights: the reader sees the aggregates, never the tables behind them.
-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'analytics_reader') THEN
        CREATE ROLE analytics_reader NOLOGIN;
    END IF;
END
$$;
-- +goose StatementEnd

GRANT USAGE ON SCHEMA analytics TO analytics_reader;

GRANT SELECT ON analytics.tenants, analytics.tenant_monthly, analytics.tenant_current TO analytics_reader;

-- +goose Down
DROP SCHEMA IF EXISTS analytics CASCADE;

DROP ROLE IF EXISTS analytics_reader;
