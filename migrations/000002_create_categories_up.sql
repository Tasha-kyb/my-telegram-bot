CREATE table categories (
user_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
id SERIAL PRIMARY KEY,
name VARCHAR(50) NOT null, 
color VARCHAR(7),
UNIQUE (user_id, name)
);