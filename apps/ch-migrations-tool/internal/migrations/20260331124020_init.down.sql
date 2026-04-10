DROP TABLE daily_outgoing_payments_mv;

--migration:split

DROP TABLE daily_outgoing_payments;

--migration:split

DROP TABLE current_payments_mv;

--migration:split

DROP TABLE current_payments;

--migration:split

DROP TABLE captured_payments_summed_mv;

--migration:split

DROP TABLE captured_payments_summed;

--migration:split

DROP TABLE payment_progress_mv;

--migration:split

DROP TABLE payment_progress;

--migration:split

DROP TABLE payment_progress_queue;
