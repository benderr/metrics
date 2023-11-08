CREATE TABLE IF NOT EXISTS metrics
(
    id text NOT NULL,
    type text NOT NULL,
    delta integer,
    value double precision,
    CONSTRAINT metrics_pkey PRIMARY KEY (id)
)