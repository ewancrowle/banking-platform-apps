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
) ENGINE = MergeTree ORDER BY (account_id, payment_id, created_at, id);

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
) ENGINE = Kafka('redpanda.redpanda.svc.cluster.local:9093', 'payments', 'clickhouse', 'JSONEachRow');

--migration:split

CREATE MATERIALIZED VIEW payments_mv TO payments AS
SELECT *
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
    account_id Int64,
    payment_id Int64,
    currency_code LowCardinality(String),
    authorised_amount AggregateFunction(sum, Int64),
    incremented_amount AggregateFunction(sum, Int64),
    is_captured AggregateFunction(max, UInt8)
)
ENGINE = AggregatingMergeTree
ORDER BY (account_id, payment_id, currency_code);

--migration:split

CREATE MATERIALIZED VIEW pending_payments_mv TO pending_payments
AS
SELECT
    account_id,
    payment_id,
    currency_code,
    sumStateIf(amount, status = 'authorised') AS authorised_amount,
    sumStateIf(amount, status = 'incremented') AS incremented_amount,
    maxState(toUInt8(status = 'captured')) AS is_captured
FROM payments
GROUP BY account_id, payment_id, currency_code;