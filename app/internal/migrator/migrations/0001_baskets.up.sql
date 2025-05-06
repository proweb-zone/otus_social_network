CREATE TABLE baskets (
    id SERIAL PRIMARY KEY,
    product_id INTEGER,
    user_id INTEGER,
    quantity INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NULL
);