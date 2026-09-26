CREATE TABLE IF NOT EXISTS users (
                                     id      BIGINT PRIMARY KEY,
                                     balance NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (balance >= 0)
    );

CREATE TABLE IF NOT EXISTS balance_history (
                                               id          BIGSERIAL PRIMARY KEY,
                                               user_id     BIGINT NOT NULL REFERENCES users(id),
    old_balance NUMERIC(12,2) NOT NULL,
    new_balance NUMERIC(12,2) NOT NULL,
    amount      NUMERIC(12,2) NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
    );

CREATE INDEX IF NOT EXISTS idx_balance_history_user_id ON balance_history(user_id);

INSERT INTO users (id, balance) VALUES (1, 500.00)
    ON CONFLICT (id) DO NOTHING;