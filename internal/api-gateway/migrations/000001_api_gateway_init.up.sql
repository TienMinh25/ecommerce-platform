-- Users table
CREATE TABLE IF NOT EXISTS users (
                                     id BIGSERIAL PRIMARY KEY,
                                     fullname VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    avatar_url VARCHAR(500),
    phone VARCHAR(15),
    birthdate DATE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    email_verified BOOLEAN DEFAULT FALSE,
    status VARCHAR(50) DEFAULT 'active',
    phone_verified BOOLEAN DEFAULT FALSE
    );

-- Roles table
CREATE TABLE IF NOT EXISTS roles (
                                     id BIGSERIAL PRIMARY KEY,
                                     role_name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    );

-- user_password table
CREATE TABLE IF NOT EXISTS user_password (
                                             id BIGSERIAL PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    );

-- addresses table
CREATE TABLE IF NOT EXISTS addresses (
                                         id BIGSERIAL PRIMARY KEY,
                                         user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    recipient_name VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    street VARCHAR(2000) NOT NULL,
    district VARCHAR(2000) NOT NULL,
    province VARCHAR(2000) NOT NULL,
    postal_code VARCHAR(20),
    country VARCHAR(2000) NOT NULL,
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    address_type VARCHAR(50) NOT NULL,
    longtitude NUMERIC(10, 7),
    latitude NUMERIC(10, 7),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    );

-- refresh token table
CREATE TABLE IF NOT EXISTS refresh_token (
    id BIGSERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    );

-- Add indexes for commonly queried fields
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_refresh_token_token ON refresh_token(token);
CREATE INDEX idx_refresh_token_user_id ON refresh_token(user_id);
CREATE INDEX idx_addresses_user_id ON addresses(user_id);

-- function trigger update timestamps when update
CREATE OR REPlACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Add triggers for updated_at timestamps
CREATE TRIGGER set_timestamp_users
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();


CREATE TRIGGER set_timestamp_roles
    BEFORE UPDATE ON roles
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();

CREATE TRIGGER set_timestamp_user_password
    BEFORE UPDATE ON user_password
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();

CREATE TRIGGER set_timestamp_addresses
    BEFORE UPDATE ON addresses
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();