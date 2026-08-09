CREATE TYPE payment_status AS ENUM('pending', 'succeeded', 'failed');


CREATE TABLE IF NOT EXISTS payments(
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(id) NOT NULL,
    amount INT CHECK (amount > 0 ) NOT NULL,
    status payment_status NOT NULL,
    method TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP(0)
)

