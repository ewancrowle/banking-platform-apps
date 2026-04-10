CREATE TABLE payment_progress_queue
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
    created_at String,
    updated_at String
) ENGINE = Kafka('redpanda.redpanda.svc.cluster.local:9093', 'payment_progress', 'clickhouse', 'JSONEachRow');

--migration:split

CREATE TABLE payment_progress
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
    created_at DateTime64(3, 'UTC'),
    updated_at DateTime64(3, 'UTC'),
) ENGINE = MergeTree ORDER BY (account_id, created_at, id);

--migration:split

CREATE MATERIALIZED VIEW payment_progress_mv TO payment_progress AS
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
    parseDateTime64BestEffort(created_at, 3, 'UTC') AS created_at,
    parseDateTime64BestEffort(updated_at, 3, 'UTC') AS updated_at
FROM payment_progress_queue;

--migration:split

CREATE TABLE captured_payments_summed
(
    account_id Int64,
    currency_code LowCardinality(String),
    amount Int64
)
ENGINE = SummingMergeTree()
ORDER BY (account_id, currency_code);

--migration:split

CREATE MATERIALIZED VIEW captured_payments_summed_mv
TO captured_payments_summed AS
SELECT
    account_id,
    currency_code,
    amount
FROM payment_progress_queue
WHERE status = 'captured';

--migration:split

CREATE TABLE current_payments
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
    created_at DateTime64(3, 'UTC'),
    updated_at DateTime64(3, 'UTC'),
) ENGINE = ReplacingMergeTree(updated_at) ORDER BY (account_id, id);

--migration:split

CREATE MATERIALIZED VIEW current_payments_mv TO current_payments AS
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
    parseDateTime64BestEffort(created_at, 3, 'UTC') AS created_at,
    parseDateTime64BestEffort(updated_at, 3, 'UTC') AS updated_at
FROM payment_progress_queue;

--migration:split

CREATE TABLE daily_outgoing_payments
(
    id Int64,
    account_id Int64,
    currency_code LowCardinality(String),
    event_date Date,
    amount Int64,
    status LowCardinality(String),
    updated_at DateTime64(3, 'UTC')
)
ENGINE = ReplacingMergeTree(updated_at)
ORDER BY (account_id, event_date, id);

--migration:split

CREATE MATERIALIZED VIEW daily_outgoing_payments_mv TO daily_outgoing_payments AS
SELECT
    id,
    account_id,
    currency_code,
    toStartOfDay(parseDateTime64BestEffort(created_at, 3, 'UTC')) AS event_date,
    amount,
    status,
    parseDateTime64BestEffort(updated_at, 3, 'UTC') AS updated_at
FROM payment_progress_queue
WHERE type IN ('withdrawal', 'card', 'account_to_account', 'fee');
