-- Write your migrate up statements here
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
);

CREATE INDEX idx_users_email ON users(email);
---- create above / drop below ----
DROP INDEX IF EXISTS idx_users_email;

DROP TABLE IF EXISTS users;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
