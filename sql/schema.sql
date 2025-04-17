BEGIN;

-- users table
-- CREATE TABLE IF NOT EXISTS users (
--     id SERIAL PRIMARY KEY,
--     username TEXT NOT NULL UNIQUE,
--     hashed_password TEXT NOT NULL,
--     created_at TIMESTAMP DEFAULT NOW()
-- );

-- devices table
CREATE TABLE IF NOT EXISTS devices (
    id SERIAL PRIMARY KEY,
    location TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
    -- user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    -- api_key TEXT NOT NULL UNIQUE,
);

-- measurements table
CREATE TABLE IF NOT EXISTS measurements (
    id SERIAL PRIMARY KEY,
    device_id INTEGER NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    value DOUBLE PRECISION NOT NULL,
    unit TEXT NOT NULL,
    timestamp TIMESTAMP DEFAULT NOW()
);

COMMIT;