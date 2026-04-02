DROP TABLE pending_payments_mv;

--migration:split

DROP TABLE pending_payments;

--migration:split

DROP TABLE total_amount_captured_mv;

--migration:split

DROP TABLE total_amount_captured;

--migration:split

DROP TABLE payments_mv;

--migration:split

DROP TABLE payments_queue;

--migration:split

DROP TABLE payments;
