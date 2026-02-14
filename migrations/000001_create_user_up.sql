CREATE table IF NOT EXISTS users (
user_id BIGINT PRIMARY key,
username VARCHAR(50) NOT null,  
created_at TIMESTAMP default CURRENT_TIMESTAMP
);