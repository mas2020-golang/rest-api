CREATE TABLE IF NOT EXISTS products
(
    id SERIAL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
    sku varchar(100),
    CONSTRAINT products_pkey PRIMARY KEY (id)
);