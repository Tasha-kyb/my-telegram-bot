CREATE table expenses (
user_id BIGINT NOT NULL references users(user_id) ON DELETE CASCADE,
category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
amount DECIMAL,
description VARCHAR(100),
created_at TIMESTAMP default CURRENT_TIMESTAMP
);

CREATE INDEX idx_expenses_user_id ON expenses(user_id);
CREATE INDEX idx_expenses_created_at ON expenses(created_at);