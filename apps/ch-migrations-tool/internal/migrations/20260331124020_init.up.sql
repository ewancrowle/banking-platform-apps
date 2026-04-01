CREATE TABLE payments
(
    id Int64,
    account_id Int64,
    payment_id Int64 DEFAULT 0,
    merchant_id Int64 DEFAULT 0,
    other_account_id Int64 DEFAULT 0,
    amount Int64,
    currency_code LowCardinality(String),
    type LowCardinality(String),
    status LowCardinality(String),
    description String,
    created_at DATETIME
) ENGINE = MergeTree ORDER BY (account_id, created_at, id);

--migration:split

CREATE TABLE payments_queue
(
    id Int64,
    account_id Int64,
    payment_id Int64,
    merchant_id Int64,
    other_account_id Int64,
    amount Int64,
    currency_code LowCardinality(String),
    type LowCardinality(String),
    status LowCardinality(String),
    description String,
    created_at DATETIME
) ENGINE = Kafka('redpanda.redpanda.svc.cluster.local:9644', 'payments', 'clickhouse', 'JSONEachRow');

--migration:split

CREATE MATERIALIZED VIEW payments_mv TO payments AS
SELECT *
FROM payments_queue;
