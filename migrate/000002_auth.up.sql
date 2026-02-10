CREATE TABLE user_history (
    user_id SERIAL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);