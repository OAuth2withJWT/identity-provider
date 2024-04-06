CREATE TYPE session_status AS ENUM ('active', 'inactive');

CREATE TABLE IF NOT EXISTS sessions (
    id SERIAL PRIMARY KEY,
    session_id VARCHAR(255) NOT NULL,
    user_id INT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    status session_status DEFAULT 'active'
);
