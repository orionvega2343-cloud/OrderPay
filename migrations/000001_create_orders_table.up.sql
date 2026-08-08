CREATE TYPE order_status AS ENUM('created', 'paid', 'cancelled', 'completed');

CREATE TABLE IF NOT EXISTS orders(
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    status order_status NOT NULL,
    total_amount INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP(0),
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP(0)
);