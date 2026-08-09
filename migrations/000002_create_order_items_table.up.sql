CREATE TABLE IF NOT EXISTS order_items(
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES orders(id),
    product_name TEXT NOT NULL,
    quantity INT CHECK (quantity > 0) NOT NULL,
    price_per_unit INT CHECK (price_per_unit > 0) NOT NULL
);