ALTER TABLE addresses
DROP COLUMN address_type;

CREATE TABLE IF NOT EXISTS address_types (
    id BIGSERIAL PRIMARY KEY,
    address_type VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER set_timestamp_address_types
BEFORE UPDATE ON address_types
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();