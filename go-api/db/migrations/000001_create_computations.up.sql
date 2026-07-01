CREATE TABLE IF NOT EXISTS computations (
    id           TEXT        PRIMARY KEY,
    username     TEXT        NOT NULL,
    input_matrix JSONB       NOT NULL,
    rows         INT         NOT NULL CHECK (rows > 0),
    cols         INT         NOT NULL CHECK (cols > 0 AND cols <= rows),
    q_matrix     JSONB       NOT NULL,
    r_matrix     JSONB       NOT NULL,
    statistics   JSONB       NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_computations_created_at ON computations (created_at DESC);
