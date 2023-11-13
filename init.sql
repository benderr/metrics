CREATE TABLE IF NOT EXISTS metrics
(
    id text NOT NULL,
    type text NOT NULL,
    delta bigint,
    value double precision,
    CONSTRAINT metrics_pkey PRIMARY KEY (id)
)