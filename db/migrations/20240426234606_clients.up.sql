CREATE TABLE IF NOT EXISTS clients (
    id SERIAL PRIMARY KEY,
    client_name VARCHAR(255) NOT NULL,
    client_id VARCHAR(255) NOT NULL UNIQUE,
    client_secret VARCHAR(255) NOT NULL,
    scope VARCHAR(255) NOT NULL,
    redirect_uri VARCHAR(255) NOT NULL  
);
 