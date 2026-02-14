CREATE table expenses (
user_id BIGINT NOT NULL references users(user_id),
amount DECIMAL,
category VARCHAR(50) NOT NULL,
description VARCHAR(100),
created_at TIMESTAMP default CURRENT_TIMESTAMP
);