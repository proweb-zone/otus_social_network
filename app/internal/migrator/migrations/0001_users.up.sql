CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(64),
    last_name VARCHAR(64),
    email VARCHAR(64),
    password VARCHAR(255),
    birth_date date,
    gender VARCHAR(12),
    hobby VARCHAR(255),
    city VARCHAR(64),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NULL
);

CREATE TABLE auth (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    token VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NULL
);

CREATE TABLE friends (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    friend_id INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
