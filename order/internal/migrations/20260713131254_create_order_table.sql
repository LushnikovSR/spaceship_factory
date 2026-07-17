-- +goose Up
CREATE TABLE IF NOT EXISTS orders (
    order_uuid       UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_uuid        UUID        NOT NULL,
    part_uuids       TEXT[]      NOT NULL DEFAULT '{}',
    total_price      NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    transaction_uuid UUID        NULL,
    payment_method   VARCHAR(32) NULL,
    status           VARCHAR(32) NOT NULL DEFAULT 'PENDING_PAYMENT'
                     CHECK (status IN ('PENDING_PAYMENT', 'PAID', 'CANCELLED'))
);

CREATE INDEX idx_orders_user_uuid ON orders (user_uuid);
CREATE INDEX idx_orders_status ON orders (status);

-- +goose Down
DROP TABLE IF EXISTS orders;