CREATE TYPE refund_status AS ENUM('approved', 'pending', 'rejected');

CREATE TABLE IF NOT EXISTS refunds(
    id SERIAL PRIMARY KEY,
    payment_id INT REFERENCES payments(id) NOT NULL,
    amount INT CHECK (amount > 0 ) NOT NULL,
    reason TEXT NOT NULL,
    status refund_status NOT NULL,
    created_at TIMESTAMPTZ DEFAULT current_timestamp(0)
)