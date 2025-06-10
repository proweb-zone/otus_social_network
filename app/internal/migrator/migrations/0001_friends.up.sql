CREATE TABLE friends (
    id SERIAL PRIMARY KEY,
    user_id INTEGER(12),
    friend_id INTEGER(12),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
);
