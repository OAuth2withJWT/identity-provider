CREATE TABLE IF NOT EXISTS verifications (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    code VARCHAR(6) NOT NULL,
    verified BOOLEAN NOT NULL
);