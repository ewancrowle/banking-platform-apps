CREATE TABLE payments
(
    id Int64,
    account_id Int64,
    merchant_id Int64,
    other_account_id Int64,
    amount Int64,
    currency_code LowCardinality(String),
    type LowCardinality(String),
    status LowCardinality(String),
    description String,
    decline_reason Int8,
    created_at DateTime64(3, 'UTC')
) ENGINE = MergeTree ORDER BY (account_id, created_at, id);

--migration:split

CREATE TABLE payments_queue
(
    id Int64,
    account_id Int64,
    merchant_id Nullable(Int64),
    other_account_id Nullable(Int64),
    amount Int64,
    currency_code LowCardinality(String),
    type LowCardinality(String),
    status LowCardinality(String),
    description String,
    decline_reason Nullable(Int8),
    created_at String
) ENGINE = Kafka('redpanda.redpanda.svc.cluster.local:9093', 'payments', 'clickhouse', 'JSONEachRow');

--migration:split

CREATE MATERIALIZED VIEW payments_mv TO payments AS
SELECT
    id,
    account_id,
    ifNull(merchant_id, 0) AS merchant_id,
    ifNull(other_account_id, 0) AS other_account_id,
    amount,
    currency_code,
    type,
    status,
    description,
    ifNull(decline_reason, 0) AS decline_reason,
    parseDateTime64BestEffort(created_at, 3, 'UTC') AS created_at
FROM payments_queue;

--migration:split

CREATE TABLE total_amount_captured
(
    account_id Int64,
    currency_code LowCardinality(String),
    total_amount Int64
)
ENGINE = SummingMergeTree
ORDER BY (account_id, currency_code);

--migration:split

CREATE MATERIALIZED VIEW total_amount_captured_mv
TO total_amount_captured AS
SELECT
    account_id,
    currency_code,
    amount AS total_amount
FROM payments
WHERE status = 'captured';

--migration:split

CREATE TABLE pending_payments
(
    id Int64,
    account_id Int64,
    currency_code LowCardinality(String),
    authorised_amount AggregateFunction(sum, Int64),
    incremented_amount AggregateFunction(sum, Int64),
    is_captured AggregateFunction(max, UInt8)
)
ENGINE = AggregatingMergeTree
ORDER BY (account_id, id, currency_code);

--migration:split

CREATE MATERIALIZED VIEW pending_payments_mv TO pending_payments
AS
SELECT
    id,
    account_id,
    currency_code,
    sumStateIf(amount, status = 'authorised') AS authorised_amount,
    sumStateIf(amount, status = 'incremented') AS incremented_amount,
    maxState(toUInt8(status = 'captured')) AS is_captured
FROM payments
GROUP BY account_id, id, currency_code;
